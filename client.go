package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/HanmaDevin/chatdev/types"
	"github.com/HanmaDevin/chatdev/views"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type Client struct {
	Conn        *websocket.Conn
	ID          string
	Chatroom    string
	Manager     *Manager
	MessageChan chan string
}

func NewClient(conn *websocket.Conn, id string, chatroom string, manager *Manager) *Client {
	return &Client{
		Conn:        conn,
		ID:          id,
		Chatroom:    chatroom,
		Manager:     manager,
		MessageChan: make(chan string),
	}
}

func (c *Client) ReadMessage(ctx echo.Context) {
	defer func() {
		c.Conn.Close()
		c.Manager.ClientEvents <- &ClientListEvent{
			Client:    c,
			EventType: "leave",
		}
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			ctx.Logger().Errorf("Error reading message: %v", err)
			return
		}
		fmt.Printf("Received: %s\n", message)

		var msgData struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(message, &msgData); err != nil {
			ctx.Logger().Errorf("Error unmarshalling message: %v", err)
			continue
		}
		c.MessageChan <- msgData.Message
	}
}

func (c *Client) WriteMessage(ctx echo.Context) {
	defer func() {
		c.Conn.Close()
		c.Manager.ClientEvents <- &ClientListEvent{
			Client:    c,
			EventType: "leave",
		}
	}()

	for text := range c.MessageChan {
		msg := types.Message{
			Sender:    c.ID,
			Message:   text,
			Timestamp: time.Now().Format("15:04"),
		}
		component := views.Msg(msg)
		var buf bytes.Buffer
		if err := component.Render(context.Background(), &buf); err != nil {
			ctx.Logger().Errorf("Error rendering message: %v", err)
			continue
		}

		fmt.Println("Writing: ", buf.String())

		for _, client := range c.Manager.Clients {
			if err := client.Conn.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
				ctx.Logger().Errorf("Error writing message: %v", err)
				return
			}
		}

	}
}
