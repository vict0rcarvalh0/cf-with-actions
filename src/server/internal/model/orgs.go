package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Orgs struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name             string         `json:"name"`
	StringConnection string         `json:"stringConnection"`
	BranchID         string         `json:"branchId"`
	PullRequestID    string         `json:"pullRequestId"`
}

func (org *Orgs) BeforeCreate(tx *gorm.DB) (err error) {
	org.ID = uuid.New()
	return
}

type CreateOrgRequest struct {
	Name             string `json:"name" validate:"required"`
	StringConnection string `json:"string_connection" validate:"required"`
	RepoName         string `json:"repo_name" validate:"required"`
	Owner            string `json:"owner"`
}

func (u *CreateOrgRequest) Validate() error {
	return validate.Struct(u)
}

func init() {
	validate = validator.New()
}
