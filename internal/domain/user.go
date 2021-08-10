package domain

import "github.com/google/uuid"

type User struct {
	ClusteredId  uint64
	Id           uuid.UUID	`json:"id"`
	EmailAddress string		`json:"emailAddress"`
	Password     string		`json:"password"`
}
