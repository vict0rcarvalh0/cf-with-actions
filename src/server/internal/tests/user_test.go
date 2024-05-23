package tests

import (
	"bytes"
	"deploy-buddy/server/internal/handler"
	"deploy-buddy/server/internal/model"
	"deploy-buddy/server/internal/repository"
	utils "deploy-buddy/server/internal/utils/slack"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

// 	// Configure a database in memory for testing
// 	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
// 	if err != nil {
// 		t.Fatalf("Erro ao abrir o banco de dados em memória: %v", err)
// 	}

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setup(t *testing.T) (*gorm.DB, *repository.UserRepository) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	path := filepath.Join(basepath, "../../.env")

	if err := godotenv.Load(path); err != nil {
		t.Fatalf("Failed to load .env file from path %s: %v", path, err)
	}

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("Failed to migrate tables: %v", err)
	}

	slackService := utils.NewSlackService()
	repo := repository.NewUserRepository(db, slackService)

	return db, repo
}

func TestUserRepository(t *testing.T) {
	db, repo := setup(t)

	// Test to Create function
	t.Run("Create", func(t *testing.T) {
		user := &model.User{Name: "John Doe", Email: "john@example.com", Password: "password123"}
		err := repo.Create(user)
		assert.NoError(t, err)

		var createdUser model.User
		err = db.First(&createdUser, "email = ?", user.Email).Error
		assert.NoError(t, err)
		assert.Equal(t, user.Name, createdUser.Name)
		assert.NotEmpty(t, createdUser.Password)
	})

// 	// Test to Delete function
// 	t.Run("Delete", func(t *testing.T) {
// 		user := &model.User{
// 			Email: "user_to_delete@example.com",
// 		}
// 		db.Create(user)
// 		err := repo.Delete(user.ID)
// 		assert.NoError(t, err)
// 	})

// 	// Test to Approve function
// 	t.Run("Approve", func(t *testing.T) {
// 		user := &model.User{
// 			Name: "User To Approve",
// 		}
// 		db.Create(user)
// 		success, err := repo.Approve(user.ID)
// 		assert.NoError(t, err)
// 		assert.True(t, success)
// 	})

// 	// Test to Decline function
// 	t.Run("Decline", func(t *testing.T) {
// 		user := &model.User{
// 			Name: "User To Decline",
// 		}
// 		db.Create(user)
// 		success, err := repo.Decline(user.ID)
// 		assert.NoError(t, err)
// 		assert.True(t, success)
// 	})
// }

// func TestUserHandler(t *testing.T) {
// 	godotenv.Load("../../.env")

// 	// Configure a database in memory for testing
// 	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
// 	assert.NoError(t, err, "Erro ao abrir o banco de dados em memória")

// 	// Migrate tables to the database
// 	err = db.AutoMigrate(&model.User{})
// 	assert.NoError(t, err, "Erro ao migrar tabelas")

// 	// Create a fake SlackService instance
// 	slackService := utils.NewSlackService()
// 	assert.NotNil(t, slackService, "SlackService não deve ser nulo")

// 	// Create User repository
// 	repo := repository.NewUserRepository(db, slackService)
// 	assert.NotNil(t, repo, "O repositório não deve ser nulo")

func setupHandler(t *testing.T) (*repository.UserRepository, *handler.UserHandler) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	path := filepath.Join(basepath, "../../.env")

	if err := godotenv.Load(path); err != nil {
		t.Fatalf("Failed to load .env file from path %s: %v", path, err)
	}

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("Failed to migrate tables: %v", err)
	}

	slackService := utils.NewSlackService()
	repo := repository.NewUserRepository(db, slackService)
	handler := handler.NewUserHandler(repo)

	return repo, handler
}

var newUser model.User

func TestUserHandler(t *testing.T) {
	_, userHandler := setupHandler(t)

	// Test to Create function
	t.Run("Create", func(t *testing.T) {
		user := &model.User{Name: "John Doe", Username: "johndoe", Email: "john.doe@email.com", Password: "password"}

		userJSON, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("Error marshalling user: %v", err)
		}

		r := chi.NewRouter()
		r.Post("/api/v1/users", userHandler.CreateUser)

		req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}

// 	// Test to GetAllUsers function
// 	t.Run("GetAllUsers", func(t *testing.T) {
// 		users := []model.User{
// 			{
// 				Name:     "User 1",
// 				Email:    "user1@example.com",
// 				Password: "password1",
// 			},
// 			{
// 				Name:     "User 2",
// 				Email:    "user2@example.com",
// 				Password: "password2",
// 			},
// 		}

// 		for _, user := range users {
// 			err := repo.Create(&user)
// 			assert.NoError(t, err, "Erro ao criar usuário")
// 		}

		// the ID is returned in the response
		// set the ID of the user to the response ID
		var createdUser model.User
		err = json.NewDecoder(rr.Body).Decode(&createdUser)
		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		newUser = createdUser

		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	// Test to GetAllUsers function
	t.Run("GetAllUsers", func(t *testing.T) {
		r := chi.NewRouter()
		r.Get("/api/v1/users", userHandler.GetAllUsers)

		req, err := http.NewRequest("GET", "/api/v1/users", nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}

		rr := httptest.NewRecorder()

// 		req := httptest.NewRequest("GET", "/users/"+user.ID.String(), nil)
// 		rr := httptest.NewRecorder()

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Test to GetUser function
	t.Run("GetUser", func(t *testing.T) {
		r := chi.NewRouter()
		r.Get("/api/v1/users/{id}", userHandler.GetUser)

		req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s", newUser.ID), nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}

		rr := httptest.NewRecorder()

// 		req := httptest.NewRequest("DELETE", "/users/"+user.ID.String(), nil)
// 		rr := httptest.NewRecorder()

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Test to UpdateUser function
	t.Run("UpdateUser", func(t *testing.T) {
		newUser.Name = "John Doe Updated"
		userJSON, err := json.Marshal(newUser)
		if err != nil {
			t.Fatalf("Error marshalling user: %v", err)
		}

		r := chi.NewRouter()
		r.Put("/api/v1/users/{id}", userHandler.UpdateUser)

		req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%s", newUser.ID), bytes.NewBuffer(userJSON))
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}

		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Test to DeleteUser function
	t.Run("DeleteUser", func(t *testing.T) {
		r := chi.NewRouter()
		r.Delete("/api/v1/users/{id}", userHandler.DeleteUser)

		req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%s", newUser.ID), nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}

		rr := httptest.NewRecorder()

// 		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
