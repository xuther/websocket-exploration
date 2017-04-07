package eventlistener

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/xuther/websocket-exploration/helpers"
)

//SaltConnection stores the token and contains the methods to interact with salt.
type SaltConnection struct {
	Token      string
	Expires    float64
	Connection *net.Conn
}

type LoginResponse struct {
	Eauth       string   `json:"eauth,omitempty"`
	Expire      float64  `json:"expire,omitempty"`
	Permissions []string `json:"perms,omitempty"`
	Start       float64  `json:"start,omitempty"`
	Token       string   `json:"token,omitempty"`
	User        string   `json:"user,omitempty"`
}

/*
	Login authenticates against the salt server and stores the connection token.

	Requires the environment variables
		- $SALT_MASTER_ADDRESS
		- $SALT_EVENT_USERNAME
		- $SALT_EVENT_PASSWORD
*/
func (sc *SaltConnection) Login() error {
	log.Printf("Logging into the salt master")

	values := make(map[string]string)
	values["username"] = os.Getenv("SALT_EVENT_USERNAME")
	values["password"] = os.Getenv("SALT_EVENT_PASSWORD")
	values["eauth"] = "pam"

	b, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", os.Getenv("SALT_MASTER_ADDRESS")+"/login", bytes.NewBuffer(b))
	if err != nil {
		log.Printf("Error building the request: %s", err.Error())
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Headers set, sending request")

	//For now ignore the certificate error, eventually we'll need to get one
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending the login request: %s", err.Error())
		return err
	}

	log.Printf("Request sent")

	respBody := make(map[string][]LoginResponse)

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading the login response: %s", err.Error())
	}

	log.Printf("Body: %s", b)
	err = json.Unmarshal(b, &respBody)
	if err != nil {
		log.Printf("Error unmarshalling login response: %s", err.Error())
	}

	log.Printf("Struct %+v", respBody)
	lr := respBody["return"][0]

	sc.Token = lr.Token
	sc.Expires = lr.Expire
	log.Printf("Done.")
	return nil
}

func (sc *SaltConnection) ListenForEvents(eventChan chan<- helpers.EventWrapper) {
	sc.Login()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", os.Getenv("SALT_MASTER_ADDRESS")+"/events", nil)
	if err != nil {
		log.Printf("Cannot open request %s", err.Error())
		return
	}

	req.Header.Add("X-Auth-Token", sc.Token)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending Request %s", err.Error())
		return
	}

	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading event" + err.Error())
		} else {
			if strings.Contains(line, "retry") {
				continue
			} else if strings.Contains(line, "tag") {
				line2, err := reader.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				}
				if strings.Contains(line2, "data") {
					log.Printf("Event Receieved")
					log.Printf("%s", line2)
					jsonString := line2[5:]
					event := helpers.EventWrapper{}

					err := json.Unmarshal([]byte(jsonString), &event)
					if err != nil {
						log.Fatal("Error unmarshalling event" + err.Error())
					}
					eventChan <- event
				}
			} else if len(line) < 1 {
				continue
			}
		}
	}
}
