package config

import (
	"deploy-buddy/server/internal/model"
	slack "deploy-buddy/server/internal/utils/slack"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Config struct {
	DB           *gorm.DB
	SlackService *slack.SlackService
}

func LoadEnv(path ...string) error {
	var err error
	if len(path) > 0 {
		err = godotenv.Load(path[0])
	} else {
		err = godotenv.Load()
	}
	return err
}

func NewConfig(path ...string) (*Config, error) {
	c := &Config{}

	if err := LoadEnv(path...); err != nil {
		return nil, err
	}

	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}

	c.DB = db
	c.DB.AutoMigrate(&model.User{}, &model.Orgs{})

	c.SlackService = slack.NewSlackService()

	return c, nil
}
