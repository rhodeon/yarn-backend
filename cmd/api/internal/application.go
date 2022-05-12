package internal

import "github.com/Mutay1/chat-backend/domain/repository"

// Application is a container to group data needed at different points throughout the server.
type Application struct {
	Config       Config
	Repositories repository.Repositories
}
