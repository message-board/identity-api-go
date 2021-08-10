package domain

import "github.com/google/uuid"

type UserRepository interface {
	CreateUser(*User)
	GetUserById(uuid.UUID) (*User, error)
	GetUserByEmailAddress(string) (*User, error)
}
