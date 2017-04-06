package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/xuther/websocket-exploration/helpers"
)

func Openws(context echo.Context) error {
	helpers.StartWebClient(context.Response().Writer, context.Request())
	return nil
}

func WriteMessage(context echo.Context) error {
	event := helpers.Event{}
	err := context.Bind(&event)
	if err != nil {
		return context.JSON(http.StatusBadRequest, "Bad event")
	}

	return helpers.WriteMessage(event)
}
