package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/HanmaDevin/chatdev/types"
	"github.com/HanmaDevin/chatdev/views"
	"github.com/labstack/echo"
)

var currentTime = time.Now().Format("15:04") + " Uhr"
var data = []types.Chat{}

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
	return nil
}

func loginAuthHandler(c echo.Context) error {
	return c.String(401, "Unauthorized")
}

func chatsHandler(c echo.Context) error {

	component := views.Chats(data)
	fmt.Printf("Data: %+v\n", data)

	return component.Render(c.Request().Context(), c.Response().Writer)
}

func convHandler(c echo.Context) error {
	id := c.Param("id")
	idint, err := strconv.Atoi(id)
	if err != nil || idint < 1 || idint > len(data) {
		return c.String(404, "Chat not found")
	}
	fmt.Printf("ID: %s\n, ID Int: %d\n, Error: %v\n", id, idint, err)
	// database simulation
	chat := data[idint-1]
	component := views.Conversation(chat)

	return component.Render(c.Request().Context(), c.Response().Writer)
}

func main() {
	e := echo.New()

	data = append(data, types.Chat{
		ID:   "1",
		Name: "Chat 1",
		Messages: []types.Message{
			{
				ID:        "1",
				Sender:    "Alice",
				Content:   "Hello, how are you?",
				Timestamp: "10:00 Uhr",
			},
			{
				ID:        "2",
				Sender:    "Bob",
				Content:   "I'm good, thanks! And you?",
				Timestamp: "10:01 Uhr",
			},
			{
				ID:        "3",
				Sender:    "Alice",
				Content:   "Doing well, just working on a project.",
				Timestamp: "10:02 Uhr",
			},
		},
	})
	data = append(data, types.Chat{
		ID:   "2",
		Name: "Chat 2",
		Messages: []types.Message{
			{
				ID:        "1",
				Sender:    "Alice",
				Content:   "Hello, how are you?",
				Timestamp: "10:00 Uhr",
			},
			{
				ID:        "2",
				Sender:    "Bob",
				Content:   "I'm good, thanks! And you?",
				Timestamp: "10:01 Uhr",
			},
			{
				ID:        "3",
				Sender:    "Alice",
				Content:   "Doing well, just working on a project.",
				Timestamp: "10:02 Uhr",
			},
		},
	})
	data = append(data, types.Chat{
		ID:   "3",
		Name: "Chat 3",
		Messages: []types.Message{
			{
				ID:        "1",
				Sender:    "Alice",
				Content:   "Hello, how are you?",
				Timestamp: "10:00 Uhr",
			},
			{
				ID:        "2",
				Sender:    "Bob",
				Content:   "I'm good, thanks! And you?",
				Timestamp: "10:01 Uhr",
			},
			{
				ID:        "3",
				Sender:    "Alice",
				Content:   "Doing well, just working on a project.",
				Timestamp: "10:02 Uhr",
			},
		},
	})

	e.GET("/", indexHandler)
	e.GET("/chats", chatsHandler)
	e.GET("/login", loginHandler)
	e.POST("/login", loginAuthHandler)
	e.GET("/chat/:id", convHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
