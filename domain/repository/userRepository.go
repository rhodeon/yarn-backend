package repository

import "github.com/Mutay1/chat-backend/models"

type UserRepository interface {
	Create(user models.User) (models.User, error)
	GetById(id string) (models.User, error)
	GetByEmail(email string) (models.User, error)
	GetByRefreshToken(refreshToken string) (models.User, error)
	UpdateRefreshToken(userId string, newRefreshToken string) error
}
