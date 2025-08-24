package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/HanmaDevin/chatdev/types"
	"github.com/HanmaDevin/chatdev/views"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	currentTime = time.Now().Format("15:04") + " Uhr"
	upgrader    = websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func indexHandler(c echo.Context) error {
	data := types.Index{
		AppName:     "ChatDev",
		CurrentTime: currentTime,
	}
	component := views.Index(data)
	fmt.Printf("Data: %+v\n", data)

	return component.Render(c.Request().Context(), c.Response().Writer)
}

func loginHandler(c echo.Context) error {
	component := views.Login(types.FormData{})
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func chatHandler(c echo.Context) error {
	component := views.Chat()
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func loginAuthHandler(c echo.Context) error {
	// Placeholder for authentication logic
	return c.Redirect(301, "/chat")
}

func registerHandler(c echo.Context) error {
	component := views.Register(types.FormData{})
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func registerPostHandler(c echo.Context) error {
	// Placeholder for registration logic
	return c.Redirect(301, "/chat")
}

func main() {
	e := echo.New()
	manager := NewManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go manager.HandleClientEvents(ctx)

	e.GET("/", indexHandler)
	e.GET("/ws", func(c echo.Context) error {
		return manager.joinChatHandler(c, ctx)
	})
	e.GET("/chat", chatHandler)
	e.GET("/login", loginHandler)
	e.POST("/login", loginAuthHandler)
	e.GET("/register", registerHandler)
	e.POST("/register", registerPostHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
