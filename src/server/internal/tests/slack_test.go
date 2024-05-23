package tests

import (
	"deploy-buddy/server/internal/model"
	utils "deploy-buddy/server/internal/utils/slack"
	"os"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func TestSlackService_SendMessage(t *testing.T) {
	godotenv.Load("../../.env")

	t.Run("send message to slack", func(t *testing.T) {
		s := utils.SlackService{
			Keys: os.Getenv("SLACK_KEYS"),
		}
		u := model.User{
			ID:         uuid.New(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			DeletedAt:  gorm.DeletedAt{},
			Name:       "name",
			Email:      "email",
			Password:   "password",
			IsApproved: true,
		}
		err := s.AskApproval(u)

		assert.Equal(t, err, nil)
	})
}
