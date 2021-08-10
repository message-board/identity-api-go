package app

import (
	"github.com/message-board/identity/internal/domain"
)

type userApp struct {
	repo domain.UserRepository
}

var _ UserAppInterface = &userApp {}

type UserAppInterface interface {
	CreateUser(*domain.User)
	// GetUserById(uuid.UUID) (*domain.User, error)
	// GetUserByEmailAddress(string) (*domain.User, error)
}

func (u *userApp) CreateUser(user *domain.User) {
	u.repo.CreateUser(user)
}