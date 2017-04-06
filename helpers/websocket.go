package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan Event
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	socket *websocket.Conn
	send   chan []byte
}

type Event struct {
	Data string `json:"data,omitempty"`
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

//manager to be used globally.
var Manager = ClientManager{
	broadcast:  make(chan Event, 100),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

//loop over to send and or receive values
func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true

		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
			}
		case event := <-manager.broadcast:
			//Marshal once
			toSend, err := json.Marshal(&event)
			if err != nil {
				continue
			}
			for conn := range manager.clients {
				select {
				case conn.send <- toSend: //this might need to be more robust - to handle a connection that is just taking a minute
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}

		}
	}
}

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func StartWebClient(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	client := &Client{socket: conn, send: make(chan []byte)}

	Manager.register <- client

	go client.write()
}

func WriteMessage(event Event) error {
	Manager.broadcast <- event
	return nil
}
