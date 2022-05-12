package repository

import "github.com/Mutay1/chat-backend/models"

type UserRepository interface {
	Create(user models.User) (models.User, error)
}
