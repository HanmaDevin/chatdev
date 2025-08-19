package types

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	upgrader = websocket.Upgrader{}
)

type Manager struct {
	ClientList          []*Client
	ClientListEventChan chan *ClientListEvent
}

type ClientListEvent struct {
	EventType string
	Client    *Client
}

func NewManager() *Manager {
	return &Manager{
		ClientList:          []*Client{},
		ClientListEventChan: make(chan *ClientListEvent),
	}
}

func (m *Manager) HandleClientListEventChan(ctx context.Context) {
	for {
		select {
		case clientListEvent, ok := <-m.ClientListEventChan:
			if !ok {
				return
			}
			switch clientListEvent.EventType {
			case "ADD":
				m.AddClient(clientListEvent.Client)
			case "REMOVE":
				m.RemoveClient(clientListEvent.Client)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) Handle(c echo.Context, ctx context.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	newClient := NewClient(ws, m)

	m.ClientListEventChan <- &ClientListEvent{
		EventType: "ADD",
		Client:    newClient,
	}

	go newClient.ReceiveMessage(c)
	go newClient.SendMessage(c, ctx)
	return nil
}

func (m *Manager) WriteMessageToChatroom(chatroom string, message string) {
	for _, client := range m.GetClientsByChatroom(chatroom) {
		client.MessageChan <- message
	}
}

func (m *Manager) AddClient(client *Client) {
	if m.GetClientById(client.Id) != nil {
		return // Client already exists, do not add again
	}
	m.ClientList = append(m.ClientList, client)
}

func (m *Manager) RemoveClient(client *Client) {
	for i, c := range m.ClientList {
		if c.Id == client.Id {
			m.ClientList = append(m.ClientList[:i], m.ClientList[i+1:]...)
			break
		}
	}
}

func (m *Manager) GetClientById(id string) *Client {
	for _, client := range m.ClientList {
		if client.Id == id {
			return client
		}
	}
	return nil
}

func (m *Manager) GetClientsByChatroom(chatroom string) []*Client {
	var clients []*Client
	for _, client := range m.ClientList {
		if client.Chatroom == chatroom {
			clients = append(clients, client)
		}
	}
	return clients
}

func (m *Manager) GetAllClients() []*Client {
	return m.ClientList
}

func (m *Manager) Clear() {
	m.ClientList = make([]*Client, 0)
}
