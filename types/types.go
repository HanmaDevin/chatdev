package types

type Index struct {
	AppName     string
	CurrentTime string
}

type Message struct {
	ID        string
	Content   string
	Timestamp string
	Sender    string
}

type ChatRoom struct {
	ID       string
	Name     string
	Messages []Message
}
