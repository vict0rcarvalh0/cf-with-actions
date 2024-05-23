package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var validate *validator.Validate

type User struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name       string         `json:"name"`
	Username   string         `json:"username"`
	Email      string         `json:"email"`
	Password   string         `json:"password"`
	IsApproved bool           `gorm:"default:false"`
	GHP        string         `json:"ghp"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=20"`
	Username string `json:"username" validate:"required,min=3,max=64"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
	GHP      string `json:"ghp"`
}	

type AuthenticateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

func (u *CreateUserRequest) Validate() error {
	return validate.Struct(u)
}

func init() {
	validate = validator.New()
}
