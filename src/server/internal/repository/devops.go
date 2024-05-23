package repository

import (
	"deploy-buddy/server/internal/model"
	github "deploy-buddy/server/internal/utils/github"
	slack "deploy-buddy/server/internal/utils/slack"
	"errors"
	"fmt"
	"log"
	"time"

	"bufio"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DevOpsRepository struct {
	db           *gorm.DB
	slackService *slack.SlackService
}

func NewDevOpsRepository(db *gorm.DB, slackService *slack.SlackService) *DevOpsRepository {
	return &DevOpsRepository{
		db:           db,
		slackService: slackService,
	}
}

func (repo *DevOpsRepository) CreateOrg(org *model.Orgs, userId uuid.UUID, owner, repoName string) (string, error) {
	log.Printf("Creating org: %s", userId)

	userRepo := NewUserRepository(repo.db, repo.slackService)
	user, err := userRepo.FindByID(userId)
	if err != nil {
		return "", errors.New("user not found")
	}

	alreadyExists := repo.db.First(&org, "string_connection = ?", org.StringConnection)
	log.Println(alreadyExists)
	if alreadyExists.Error == nil {
		return "", errors.New("organization already exists")
	}

	org.ID = uuid.New()

	if owner == "" {
		owner = user.Username
	}

	go repo.RetrieveMetadatas(user, org, user.Username, user.GHP, owner, repoName)

	return "Organization created successfully", nil
}

func (repo *DevOpsRepository) RetrieveMetadatas(user *model.User, org *model.Orgs, githubUsername, githubToken, owner, repoName string) (string, error) {

	gc := github.NewGithubClient(githubUsername, githubToken)

	log.Println("Retrieving metadata")
	folderName := fmt.Sprintf("./metadatas/%s", org.ID.String())

	file, err := os.Create("temp.txt")
	if err != nil {
		log.Println("Error creating file:", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(org.StringConnection)
	if err != nil {
		log.Println("Error writing to file:", err)
	}

	err = writer.Flush()
	if err != nil {
		log.Println("Error flushing file:", err)
	}

	cmdLogin := exec.Command("sfdx", "force:auth:sfdxurl:store", "-f", "temp.txt", "-a", org.ID.String())

	cmd := exec.Command("sfdx", "force:mdapi:retrieve", "-r", folderName, "-o", org.ID.String(), "-k", "./metadatas/package.xml")

	cmdUnpack := exec.Command("unzip", folderName+"/unpackaged.zip", "-d", folderName)

	cmdDeleteZip := exec.Command("rm", "-rf", folderName+"/unpackaged.zip")

	errLogin := cmdLogin.Run()
	err = cmd.Run()
	errUnpack := cmdUnpack.Run()
	errDelete := cmdDeleteZip.Run()

	exec.Command("mv", folderName+"/unpackaged", folderName).Run()

	if err != nil || errUnpack != nil || errDelete != nil || errLogin != nil {
		log.Println("Error:", err)
		return "", err
	}

	err = os.Remove("temp.txt")
	if err != nil {
		log.Println("Error removing file:", err)
	}

	log.Println("Arquivo deletado com sucesso.")

	branchName := org.ID.String()
	branchID, err := gc.CreateBranch(owner, repoName, "prod", branchName)
	if err != nil {
		return "", fmt.Errorf("failed to create branch: %v", err)
	}
	org.BranchID = branchID

	if err := gc.AddFilesToRepo(owner, repoName, folderName, branchName); err != nil {
		return "", fmt.Errorf("failed to add files: %v", err)
	}

	prTitle := fmt.Sprintf("[%s] - Deploy metadata update %s", org.ID.String(), time.Now().Format("2006-01-02 15:04:05"))

	prID, err := gc.OpenPullRequest(user, owner, repoName, prTitle, branchName, "prod")
	if err != nil {
		return "", fmt.Errorf("failed to open pull request: %v", err)
	}

	log.Println("prID:", prID)

	org.PullRequestID = prID

	repo.db.Create(org)
	return "Metadata deployed and pull request opened successfully", nil
}
