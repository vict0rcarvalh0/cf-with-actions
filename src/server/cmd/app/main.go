package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"deploy-buddy/server/internal/config"
	"deploy-buddy/server/internal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	c, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	port := os.Getenv("PORT")
	log.Printf("Port: %v\n", port)

	r := chi.NewRouter()
	ur := routes.NewUsersRoutes(c)
	sr := routes.NewSlackRoutes(c)
	dr := routes.NewDevOpsRoutes(c)

	r.Use(middleware.Logger)
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", ur.Router)
		r.Route("/slack", sr.Router)
		r.Route("/devops", dr.Router)
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// { "message": "Welcome to Deploy Buddy API", "status": 200 }
		render.JSON(w, r, map[string]string{"message": "Welcome to Deploy Buddy API", "status": "200"})
	})

	host := "http://localhost"
	log.Printf("Server running at %v:%v\n", host, port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
