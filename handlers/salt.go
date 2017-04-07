package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/xuther/websocket-exploration/eventlistener"
	"github.com/xuther/websocket-exploration/helpers"
)

func Login(context echo.Context) error {

	sc := &eventlistener.SaltConnection{}
	err := sc.Login()
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}
	return context.JSON(http.StatusOK, sc)
}

func Listen(context echo.Context) error {

	sc := &eventlistener.SaltConnection{}
	go sc.ListenForEvents(helpers.Manager.Broadcast)
	return nil
}
