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

type Chat struct {
	ID       string
	Name     string
	Messages []Message
}
