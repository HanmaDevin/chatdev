package types

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/HanmaDevin/chatdev/templates"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	pingInterval = 9 * time.Second
	pongWait     = 10 * time.Second
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

func (c *Client) SendMessage(echoContext echo.Context, ctx context.Context) {
	defer func() {
		c.Connection.Close()
		c.Manager.ClientListEventChan <- &ClientListEvent{
			EventType: "REMOVE",
			Client:    c,
		}
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case msg, ok := <-c.MessageChan:
			if !ok {
				return
			}
			component := templates.Message(msg)
			buffer := new(bytes.Buffer)
			if err := component.Render(ctx, buffer); err != nil {
				echoContext.Logger().Errorf("Error rendering message: %v", err)
				return
			}

			err := c.Connection.WriteMessage(websocket.TextMessage, buffer.Bytes())
			if err != nil {
				echoContext.Logger().Errorf("Error sending message to client %s: %v", c.Id, err)
				continue
			}

		case <-ctx.Done():
			return

		case <-ticker.C:
			if err := c.Connection.WriteMessage(websocket.PingMessage, []byte("")); err != nil {
				echoContext.Logger().Errorf("Error sending ping: %v", err)
				return
			}
			fmt.Printf("Ping sent to client %s\n", c.Id)
		}
	}
}

func (c *Client) ReceiveMessage(ctx echo.Context) {
	if err := c.Connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		ctx.Logger().Errorf("Error setting read deadline: %v", err)
		return
	}

	c.Connection.SetPongHandler(func(appData string) error {
		if err := c.Connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			ctx.Logger().Errorf("Error resetting read deadline: %v", err)
			return err
		}
		fmt.Printf("Pong received from client %s\n", c.Id)
		return nil
	})
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
		c.Manager.WriteMessageToChatroom(c.Chatroom, string(msg))
	}
}
