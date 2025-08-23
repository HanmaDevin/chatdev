package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/labstack/echo"
)

type Manager struct {
	Clients      []*Client
	ClientEvents chan *ClientListEvent
}

type ClientListEvent struct {
	Client    *Client
	EventType string // "join" or "leave"
}

func NewManager() *Manager {
	return &Manager{
		Clients:      make([]*Client, 0),
		ClientEvents: make(chan *ClientListEvent),
	}
}

func generateRandomName() string {
	adjectives := []string{"Swift", "Brave", "Clever", "Mighty", "Silent", "Bright", "Lucky", "Fierce", "Gentle", "Bold"}
	nouns := []string{"Lion", "Falcon", "Tiger", "Wolf", "Eagle", "Panther", "Fox", "Bear", "Shark", "Dragon"}

	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	number := rand.Intn(1000)

	return adj + noun + fmt.Sprintf("%03d", number)
}

func (m *Manager) HandleClientEvents(ctx context.Context) {
	for event := range m.ClientEvents {
		select {
		case <-ctx.Done():
			return
		default:
			switch event.EventType {
			case "join":
				for _, client := range m.Clients {
					if client.ID == event.Client.ID {
						return
					}
				}
				m.Clients = append(m.Clients, event.Client)
			case "leave":
				for i, client := range m.Clients {
					if client.ID == event.Client.ID {
						m.Clients = append(m.Clients[:i], m.Clients[i+1:]...)
						break
					}
				}
			}
		}
	}
}

func (m *Manager) joinChatHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	newClient := NewClient(ws, generateRandomName(), "default", m)

	m.ClientEvents <- &ClientListEvent{
		Client:    newClient,
		EventType: "join",
	}

	go newClient.ReadMessage(c)
	go newClient.WriteMessage(c)

	return nil
}
