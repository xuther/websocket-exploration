package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/xuther/websocket-exploration/handlers"
	"github.com/xuther/websocket-exploration/helpers"
)

func main() {
	port := ":42042"
	router := echo.New()

	router.GET("/ws", handlers.Openws)
	router.POST("/event", handlers.WriteMessage)
	router.GET("/login", handlers.Login)
	router.GET("/listen", handlers.Listen)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}
	go helpers.Manager.Start()

	router.StartServer(&server)

}
