package dto

import "time"

type MessageDTO struct {
	ID        string
	Message   string
	From      string
	Timestamp time.Time
}
