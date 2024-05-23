package routes

import (
	"deploy-buddy/server/internal/config"
	"deploy-buddy/server/internal/handler"
	middlewares "deploy-buddy/server/internal/middleware/jwt"
	"deploy-buddy/server/internal/repository"

	"github.com/go-chi/chi/v5"
)

type UsersRoutes struct {
	Handler *handler.UserHandler
	config  *config.Config
	repo    *repository.UserRepository
}

func NewUsersRoutes(config *config.Config) UsersRoutes {
	repo := repository.NewUserRepository(config.DB, config.SlackService)
	handler := handler.NewUserHandler(repo)

	return UsersRoutes{
		Handler: handler,
		config:  config,
		repo:    repo,
	}
}

func (ur *UsersRoutes) Router(r chi.Router) {
	r.Post("/", ur.Handler.CreateUser)
	r.Post("/auth", ur.Handler.AuthenticateUser)
	r.Group(func(r chi.Router) {
		r.Use(middlewares.JWTAuthMiddleware)
		r.Get("/", ur.Handler.GetAllUsers)
		r.Get("/{id}", ur.Handler.GetUser)
		r.Delete("/{id}", ur.Handler.DeleteUser)
		r.Put("/{id}", ur.Handler.UpdateUser)
	})
	r.Post("/approve/{id}", ur.Handler.ApproveUser)
	r.Post("/decline/{id}", ur.Handler.DeclineUser)
}
