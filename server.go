package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/xuther/websocket-exploration/handlers"
)

func main() {
	port := ":42042"
	router := echo.New()

	router.GET("/ws", handlers.Openws)
	router.POST("/event", handlers.WriteMessage)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}
	router.StartServer(&server)
}
