package types

type User struct {
	ID       string
	Username string
	Password string
	Chats    []Chat
}
