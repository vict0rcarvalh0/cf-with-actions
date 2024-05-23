package github

import (
	"context"
	"deploy-buddy/server/internal/model"
	slack "deploy-buddy/server/internal/utils/slack"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
)

type GithubClient struct {
	githubClient *github.Client
	slack        *slack.SlackService
}

func NewGithubClient(username, token string) *GithubClient {
	log.Printf("Initializing GitHub client with username: %s and token: %s", username, token)

	tp := github.BasicAuthTransport{
		Username: username,
		Password: token,
	}
	client := github.NewClient(tp.Client())

	s := slack.NewSlackService()

	return &GithubClient{
		githubClient: client,
		slack:        s,
	}
}

var ctx = context.Background()

func (gc *GithubClient) AddFilesToRepo(owner, repoName, folderPath, branchName string) error {
	log.Printf("Adding files to repo: %s/%s", owner, repoName)
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error walking the path %v: %v", path, err)
			return err
		}
		if !info.IsDir() {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Error reading the file %v: %v", path, err)
				return err
			}
			githubPath := path[len(folderPath)+1:]
			return gc.CreateOrUpdateFile(owner, repoName, githubPath, "Add/update file", fileContent, branchName)
		}
		return nil
	})

	log.Println("Files added to repo")
	return err
}

func (gc *GithubClient) CreateOrUpdateFile(owner, repoName, path, message string, content []byte, branch string) error {
	log.Printf("Creating or updating file: %s", path)
	opts := &github.RepositoryContentGetOptions{Ref: branch}
	fileContent, _, _, err := gc.githubClient.Repositories.GetContents(ctx, owner, repoName, path, opts)
	if err != nil {
		if _, ok := err.(*github.ErrorResponse); ok {
			log.Println("File does not exist, creating it")
			createOptions := &github.RepositoryContentFileOptions{
				Message: github.String(message),
				Content: content,
				Branch:  github.String(branch),
			}
			_, _, err = gc.githubClient.Repositories.CreateFile(ctx, owner, repoName, path, createOptions)
			return err
		}
		log.Printf("Error getting contents of the file %v: %v", path, err)
		return err
	}

	log.Println("File exists, updating it")
	updateOptions := &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: content,
		SHA:     fileContent.SHA,
		Branch:  github.String(branch),
	}
	_, _, err = gc.githubClient.Repositories.UpdateFile(ctx, owner, repoName, path, updateOptions)

	log.Println("File updated")
	return err
}

func (gc *GithubClient) CreateBranch(owner, repoName, baseBranch, newBranch string) (string, error) {
	log.Printf("Attempting to create branch '%s' from '%s' in repo '%s/%s'", newBranch, baseBranch, owner, repoName)

	refs, _, err := gc.githubClient.Git.ListRefs(ctx, owner, repoName, nil)
	if err != nil {
		log.Printf("Error listing refs for repo '%s/%s': %v", owner, repoName, err)
		return "", err
	}

	for _, ref := range refs {
		if ref.GetRef() == "refs/heads/"+newBranch {
			log.Printf("Branch '%s' already exists in repo '%s/%s'", newBranch, owner, repoName)
			return "", fmt.Errorf("branch '%s' already exists", newBranch)
		}
	}

	ref, _, err := gc.githubClient.Git.GetRef(ctx, owner, repoName, "refs/heads/"+baseBranch)
	if err != nil {
		log.Printf("Failed to get base branch ref '%s' from repo '%s/%s': %v", baseBranch, owner, repoName, err)
		return "", err
	}

	newRef := &github.Reference{Ref: github.String("refs/heads/" + newBranch), Object: &github.GitObject{SHA: ref.Object.SHA}}
	gitRef, _, err := gc.githubClient.Git.CreateRef(ctx, owner, repoName, newRef)
	if err != nil {
		log.Printf("Failed to create branch '%s' in repo '%s/%s': %v", newBranch, owner, repoName, err)
		return "", err
	}

	log.Printf("Branch '%s' created successfully in repo '%s/%s'", newBranch, owner, repoName)
	return gitRef.GetRef(), nil
}

func (gc *GithubClient) OpenPullRequest(user *model.User, owner, repoName, title, head, base string) (string, error) {
	log.Printf("Opening pull request: %s", title)
	newPR := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(head),
		Base:  github.String(base),
	}
	prRef, _, err := gc.githubClient.PullRequests.Create(ctx, owner, repoName, newPR)
	if err != nil {
		log.Printf("Failed to create pull request: %v", err)
		return "", err
	}

	prURL := prRef.GetHTMLURL()

	// go gc.slack.NotifyPullRequestCreated(user, &title, &repoName, &head, &prURL)
	err = gc.slack.NotifyPullRequestCreated(user, &title, &repoName, &head, &prURL)
	if err != nil {
		log.Printf("Failed to send notification to Slack: %v", err)
		return "", err
	}

	log.Printf("Notification sent to Slack for pull request: %s - %s", title, prURL)
	return prURL, nil
}
