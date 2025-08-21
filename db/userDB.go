package db

import (
	"fmt"

	"github.com/HanmaDevin/chatdev/types"
)

var users = []types.User{
	{
		ID:       "1",
		Username: "admin",
		Password: "admin123",
		Chats:    []types.Chat{},
	},
}

func SaveUser(user types.User) error {
	// Simulate saving to a database
	// In a real application, you would use a database driver to save the user
	// For example, using GORM or sqlx for SQL databases
	// Or using a NoSQL database driver like MongoDB or Redis

	// Here we just print the user to simulate saving
	fmt.Printf("Saving user: %+v\n", user)
	users = append(users, user)
	return nil
}

func ValidUser(username, password string) (*types.User, error) {
	for _, user := range users {
		if user.Username == username && user.Password == password {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("invalid username or password")
}

func GetUserByID(id int) (*types.User, error) {
	if id < 1 || id > len(users) {
		return nil, fmt.Errorf("user with ID %d not found", id)
	}
	user := users[id-1]
	return &user, nil
}

func GetAllUsers() []types.User {
	// Return a copy of the users slice to avoid external modifications
	usersCopy := make([]types.User, len(users))
	copy(usersCopy, users)
	return usersCopy
}

func GetUserByUsername(username string) (*types.User, error) {
	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with username %s not found", username)
}
