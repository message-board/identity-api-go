package events

import "github.com/google/uuid"

type UserCreated struct {
	Id           uuid.UUID `json:"id"`
	EmailAddress string    `json:"emailAddress"`
}
