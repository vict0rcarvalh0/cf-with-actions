package repository

import (
	"deploy-buddy/server/internal/model"
	utils "deploy-buddy/server/internal/utils/slack"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

// UserRepository representa o repositório de usuários.
type UserRepository struct {
	db           *gorm.DB
	slackService *utils.SlackService
}

// NewUserRepository cria um novo repositório de usuários.
func NewUserRepository(db *gorm.DB, slackService *utils.SlackService) *UserRepository {
	return &UserRepository{db: db, slackService: slackService}
}

// Create cria um novo usuário no banco de dados.
func (repo *UserRepository) Create(user *model.User) error {
	var existingUser model.User
	if err := repo.db.Where("email = ? OR username = ?", user.Email, user.Username).First(&existingUser).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to query user: %v", err)
		}
	} else {
		return fmt.Errorf("user with the same email or username already exists")
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(passHash)

	if err := repo.db.Create(user).Error; err != nil {
		return err
	}

	go repo.slackService.AskApproval(*user)

	return nil
}

func (repo *UserRepository) Authenticate(email, password string) (string, error) {
	var user model.User
	err := repo.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return "", fmt.Errorf("user not found: %v", err)
	}

	log.Printf("User found: %v\n", user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid password: %v", err)
	}

	log.Printf("User authenticated: %v", user)

	//Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"id":    user.ID,
		"name":  user.Name,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	log.Printf("Token: %v\n", token)

	generatedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	log.Printf("Generated token: %v\n", generatedToken)

	return generatedToken, nil
}

// FindAll retorna todos os usuários.
func (repo *UserRepository) FindAll() ([]model.User, error) {
	var users []model.User
	err := repo.db.Find(&users).Error
	return users, err
}

// FindByID encontra um usuário pelo UUID.
func (repo *UserRepository) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := repo.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

// Update atualiza um usuário.
func (repo *UserRepository) Update(user *model.User, passwordChanged bool) error {
	if passwordChanged {
		passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user.Password = string(passHash)
	}

	return repo.db.Save(user).Error
}

// Delete deleta um usuário.
func (repo *UserRepository) Delete(id uuid.UUID) error {
	return repo.db.Delete(&model.User{}, id).Error
}

func (repo *UserRepository) Approve(id uuid.UUID) (bool, error) {
	user, err := repo.FindByID(id)
	if err != nil {
		return false, err
	}

	user.IsApproved = true

	err = repo.Update(user, false)
	if err != nil {
		return false, err
	}

	go repo.slackService.Approved(*user)

	return true, nil
}

func (repo *UserRepository) Decline(id uuid.UUID) (bool, error) {
	user, err := repo.FindByID(id)
	if err != nil {
		return false, err
	}

	if user.IsApproved {
		return false, errors.New("cannot decline a user that has been approved")
	}

	err = repo.Delete(user.ID)
	if err != nil {
		return false, err
	}

	go repo.slackService.Declined(*user)

	return true, nil
}
