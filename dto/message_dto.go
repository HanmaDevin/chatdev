package dto

import "time"

type MessageDTO struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	From      string    `json:"from"`
	Timestamp time.Time `json:"timestamp"`
}
