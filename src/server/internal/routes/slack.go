package routes

import (
	"deploy-buddy/server/internal/config"
	"deploy-buddy/server/internal/handler"
	"deploy-buddy/server/internal/repository"

	"github.com/go-chi/chi/v5"
)

type SlackRoutes struct {
	config  *config.Config
	Handler *handler.SlackHandler
}

func NewSlackRoutes(config *config.Config) SlackRoutes {
	ur := repository.NewUserRepository(config.DB, config.SlackService)
	handler := handler.NewSlackHandler(ur)

	return SlackRoutes{
		config:  config,
		Handler: handler,
	}

}

func (sr *SlackRoutes) Router(r chi.Router) {
	r.Post("/interactive", sr.Handler.Interactive)
}
