package domain

import "github.com/google/uuid"

type User struct {
	Id           uuid.UUID
	EmailAddress string
	Password     string
}
