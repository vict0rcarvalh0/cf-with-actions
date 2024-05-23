package handler

import (
	"encoding/json"
	"net/http"
	"time"

	middlewares "deploy-buddy/server/internal/middleware/jwt"
	"deploy-buddy/server/internal/model"
	"deploy-buddy/server/internal/repository"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type DevOpsHandler struct {
	repo *repository.DevOpsRepository
}

func NewDevOpsHandler(repo *repository.DevOpsRepository) *DevOpsHandler {
	return &DevOpsHandler{repo: repo}
}

func (h *DevOpsHandler) CreateOrg(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()

	var createOrgRequest model.CreateOrgRequest
	err := json.NewDecoder(r.Body).Decode(&createOrgRequest)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request payload", "message": err.Error()})
		return
	}

	userID, err := uuid.Parse(r.Context().Value(middlewares.UserContextKey).(string))
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "Unauthorized", "message": "Invalid user ID"})
		return
	}

	if err := createOrgRequest.Validate(); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	org := &model.Orgs{
		Name:             createOrgRequest.Name,
		StringConnection: createOrgRequest.StringConnection,
	}

	message, err := h.repo.CreateOrg(org, userID, createOrgRequest.Owner, createOrgRequest.RepoName)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to create user", "message": err.Error()})
		return
	}

	t1 := time.Now()

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"message":      message,
		"elapsed_time": t1.Sub(t0).String(),
	})
}
