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

var (
	pongWaitTime = 10 * time.Second
	pingInterval = 9 * time.Second
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
	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWaitTime)); err != nil {
		ctx.Logger().Errorf("Error setting read deadline: %v", err)
		return
	}
	c.Conn.SetPongHandler(func(string) error {
		ctx.Logger().Info("Received pong")
		fmt.Println("Pong received from client")
		c.Conn.SetReadDeadline(time.Now().Add(pongWaitTime))
		return nil
	})
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
		if err := c.Manager.WriteMessage(msgData.Message, "default"); err != nil {
			ctx.Logger().Errorf("Error writing message: %v", err)
			return
		}
	}
}

func (c *Client) WriteMessage(echoCtx echo.Context, ctx context.Context) {
	defer func() {
		c.Conn.Close()
		c.Manager.ClientEvents <- &ClientListEvent{
			Client:    c,
			EventType: "leave",
		}
	}()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()
	for {
		select {
		case text, ok := <-c.MessageChan:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			msg := types.Message{
				Sender:    c.ID,
				Message:   text,
				Timestamp: time.Now().Format("15:04"),
			}
			component := views.Msg(msg)
			var buf bytes.Buffer
			if err := component.Render(ctx, &buf); err != nil {
				echoCtx.Logger().Errorf("Error rendering message: %v", err)
				continue
			}

			fmt.Println("Writing: ", buf.String())
			if err := c.Conn.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
				echoCtx.Logger().Errorf("Error writing message: %v", err)
				return
			}
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte("Ping!")); err != nil {
				return
			}
			fmt.Println("Ping sent to client")
		}
	}
}
