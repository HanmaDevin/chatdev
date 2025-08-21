package db

import (
	"fmt"

	"github.com/HanmaDevin/chatdev/types"
)

var chats = []types.Chat{
	{
		ID:       "1",
		Name:     "General",
		Messages: []types.Message{},
	},
}

func SaveChat(chat types.Chat) error {
	// Simulate saving to a database
	// In a real application, you would use a database driver to save the chat
	// For example, using GORM or sqlx for SQL databases
	// Or using a NoSQL database driver like MongoDB or Redis

	// Here we just print the chat to simulate saving
	fmt.Printf("Saving chat: %+v\n", chat)
	chats = append(chats, chat)
	return nil
}

func GetChatByID(id int) (*types.Chat, error) {
	if id < 1 || id > len(chats) {
		return nil, fmt.Errorf("chat with ID %d not found", id)
	}
	chat := chats[id-1]
	return &chat, nil
}

func GetAllChats() []types.Chat {
	// Return a copy of the chats slice to avoid external modifications
	chatsCopy := make([]types.Chat, len(chats))
	copy(chatsCopy, chats)
	return chatsCopy
}
