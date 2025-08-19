package main

import (
	"context"

	"github.com/HanmaDevin/chatdev/templates"
	"github.com/HanmaDevin/chatdev/types"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	manager := types.NewManager()

	go manager.HandleClientListEventChan(context.Background())

	e.GET("/", func(c echo.Context) error {
		component := templates.Index()
		return component.Render(context.Background(), c.Response().Writer)
	})
	e.GET("/chatroom", manager.Handle)

	e.Logger.Fatal(e.Start(":8080"))
}
