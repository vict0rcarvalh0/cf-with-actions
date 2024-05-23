package handler

import (
	"encoding/json"
	"net/http"

	"deploy-buddy/server/internal/model"
	"deploy-buddy/server/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var createUserRequest model.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request payload", "message": err.Error()})
		return
	}

	if err := createUserRequest.Validate(); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	user := &model.User{
		Name:     createUserRequest.Name,
		Username: createUserRequest.Username,
		Email:    createUserRequest.Email,
		Password: createUserRequest.Password,
		GHP:      createUserRequest.GHP,
	}

	err = h.repo.Create(user)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to create user", "message": err.Error()})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

func (h *UserHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var authenticateUserRequest model.AuthenticateUserRequest

	err := json.NewDecoder(r.Body).Decode(&authenticateUserRequest)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request payload", "message": err.Error()})
		return
	}

	user, err := h.repo.Authenticate(authenticateUserRequest.Email, authenticateUserRequest.Password)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "Invalid credentials", "message": err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.FindAll()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to fetch users", "message": err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	var userID uuid.UUID
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid user ID", "message": err.Error()})
		return
	}

	user, err := h.repo.FindByID(userID)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "User not found", "message": err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var userID uuid.UUID
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid user ID", "message": err.Error()})
		return
	}

	var user model.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request payload", "message": err.Error()})
		return
	}

	user.ID = userID

	updatingPassword := user.Password != ""
	err = h.repo.Update(&user, updatingPassword)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to update user", "message": err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var userID uuid.UUID
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid user ID", "message": err.Error()})
		return
	}

	err = h.repo.Delete(userID)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to delete user", "message": err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "User deleted"})
}

func (h *UserHandler) ApproveUser(w http.ResponseWriter, r *http.Request) {
	var userID uuid.UUID
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid user ID", "message": err.Error()})
		return
	}

	approved, err := h.repo.Approve(userID)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to approve user", "message": err.Error()})
		return
	}

	if !approved {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "User not found"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "User approved"})
}

func (h *UserHandler) DeclineUser(w http.ResponseWriter, r *http.Request) {
	var userID uuid.UUID
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid user ID", "message": err.Error()})
		return
	}

	declined, err := h.repo.Decline(userID)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to decline user", "message": err.Error()})
		return
	}

	if !declined {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "User not found"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "User declined"})
}
