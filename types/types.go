package types

type Index struct {
	AppName     string
	CurrentTime string
}

type Message struct {
	Content   string
	Timestamp string
}

type Chat struct {
	Messages []Message
}
