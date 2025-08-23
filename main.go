package main

import (
	"bytes"
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

func joinChatHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// _, msg, err := ws.ReadMessage()
		// if err != nil {
		// 	break
		// }
		// fmt.Printf("Received: %s\n", msg)
		currentTime = time.Now().Format("15:04") + " Uhr"
		msg := types.Message{
			Sender:    "System",
			Message:   "Hello! This is a test message. To see how longer messages are handled, please wait for the next message.",
			Timestamp: currentTime,
		}

		component := views.Msg(msg)
		buffer := new(bytes.Buffer)
		component.Render(context.Background(), buffer)

		if err := ws.WriteMessage(websocket.TextMessage, buffer.Bytes()); err != nil {
			c.Logger().Error(err)
		}
		time.Sleep(3 * time.Second)
	}
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

	e.GET("/", indexHandler)
	e.GET("/ws", joinChatHandler)
	e.GET("/chat", chatHandler)
	e.GET("/login", loginHandler)
	e.POST("/login", loginAuthHandler)
	e.GET("/register", registerHandler)
	e.POST("/register", registerPostHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
