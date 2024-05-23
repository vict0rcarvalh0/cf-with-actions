package routes

import (
	"deploy-buddy/server/internal/config"
	"deploy-buddy/server/internal/handler"
	middlewares "deploy-buddy/server/internal/middleware/jwt"
	"deploy-buddy/server/internal/repository"

	"github.com/go-chi/chi/v5"
)

type DevOpsRoutes struct {
	config  *config.Config
	repo    *repository.DevOpsRepository
	Handler *handler.DevOpsHandler
}

func NewDevOpsRoutes(config *config.Config) DevOpsRoutes {
	repo := repository.NewDevOpsRepository(config.DB, config.SlackService)
	handler := handler.NewDevOpsHandler(repo)

	return DevOpsRoutes{
		Handler: handler,
		config:  config,
		repo:    repo,
	}
}

func (ur *DevOpsRoutes) Router(r chi.Router) {
	r.Use(middlewares.JWTAuthMiddleware)
	r.Post("/", ur.Handler.CreateOrg)
}
