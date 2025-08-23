package types

type Index struct {
	AppName     string
	CurrentTime string
}

type Message struct {
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type Chat struct {
	Messages []Message
}
