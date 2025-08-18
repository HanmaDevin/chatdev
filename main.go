package main

import (
	"context"

	"github.com/HanmaDevin/chatdev/templates"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	upgrader = websocket.Upgrader{}
)

func hello(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	for {
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, World!"))
		if err != nil {
			c.Logger().Error("write message error:", err)
		}

		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error("read message error:", err)
			return err
		}
		c.Logger().Info("Received message:", string(msg))
	}
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		component := templates.Index()
		return component.Render(context.Background(), c.Response().Writer)
	})
	e.GET("/chatroom", hello)

	e.Logger.Fatal(e.Start(":8080"))
}
