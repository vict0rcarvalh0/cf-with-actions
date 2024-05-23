package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"deploy-buddy/server/internal/repository"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type SlackHandler struct {
	repo *repository.UserRepository
}

func NewSlackHandler(repo *repository.UserRepository) *SlackHandler {
	return &SlackHandler{repo: repo}
}

func (h *SlackHandler) Interactive(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	payload := r.FormValue("payload")

	type InteractionPayload struct {
		Type string `json:"type"`
		User struct {
			ID string `json:"id"`
		} `json:"user"`
		Actions []struct {
			Name     string `json:"name"`
			Value    string `json:"value"`
			Type     string `json:"type"`
			ActionID string `json:"action_id"`
		} `json:"actions"`
		ResponseURL string `json:"response_url"`
	}

	var interaction InteractionPayload
	if err := json.Unmarshal([]byte(payload), &interaction); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid interaction payload"})
		return
	}

	log.Printf("Interaction: %v\n", interaction)

	for _, action := range interaction.Actions {
		parts := strings.SplitN(action.Value, "_", 2)
		if len(parts) != 2 {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid action value"})
			return
		}
		actionType := parts[0]
		userIDString := parts[1]

		userID, err := uuid.Parse(userIDString)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid user ID format"})
			return
		}

		if actionType == "approve" {
			_, err = h.repo.Approve(userID)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, map[string]string{"error": "Failed to approve user", "message": err.Error()})
				return
			}
		} else if actionType == "decline" {
			_, err = h.repo.Decline(userID)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, map[string]string{"error": "Failed to decline user", "message": err.Error()})
				return
			}
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Interaction processed successfully"})
}
