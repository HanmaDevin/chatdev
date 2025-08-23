package db

import (
	"github.com/HanmaDevin/chatdev/types"
	"golang.org/x/crypto/bcrypt"
)

var hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("admin123"), 8)

var users = []types.User{
	{
		ID:       "1",
		Username: "admin",
		Password: string(hashedPassword),
		Chats:    []types.Chat{},
	},
}
