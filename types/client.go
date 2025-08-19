package types

import (
	"bytes"
	"context"
	"fmt"

	"github.com/HanmaDevin/chatdev/templates"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type Client struct {
	Id          string
	Connection  *websocket.Conn
	Chatroom    string
	Manager     *Manager
	MessageChan chan string
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		Id:          uuid.New().String(),
		Connection:  conn,
		Chatroom:    "general",
		Manager:     manager,
		MessageChan: make(chan string),
	}
}

func (c *Client) SendMessage(ctx echo.Context) {
	defer func() {
		c.Connection.Close()
		c.Manager.ClientListEventChan <- &ClientListEvent{
			EventType: "REMOVE",
			Client:    c,
		}
	}()

	for {
		select {
		case msg, ok := <-c.MessageChan:
			if !ok {
				return
			}
			component := templates.Message(msg)
			buffer := new(bytes.Buffer)
			if err := component.Render(context.Background(), buffer); err != nil {
				ctx.Logger().Errorf("Error rendering message: %v", err)
				return
			}

			for _, client := range c.Manager.ClientList {
				err := client.Connection.WriteMessage(websocket.TextMessage, buffer.Bytes())
				if err != nil {
					ctx.Logger().Errorf("Error sending message to client %s: %v", client.Id, err)
					continue
				}
			}

		case <-context.Background().Done():
			return
		}
	}
}

func (c *Client) ReceiveMessage(ctx echo.Context) {
	defer func() {
		c.Connection.Close()
		c.Manager.ClientListEventChan <- &ClientListEvent{
			EventType: "REMOVE",
			Client:    c,
		}
	}()
	for {
		_, msg, err := c.Connection.ReadMessage()
		if err != nil {
			ctx.Logger().Errorf("Error reading message: %v", err)
			return
		}

		fmt.Printf("Received message: %s\n", msg)
		c.MessageChan <- string(msg)
	}
}
