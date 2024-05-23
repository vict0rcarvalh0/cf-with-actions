package utils

import (
	"deploy-buddy/server/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type SlackService struct {
	Keys string
}

func NewSlackService() *SlackService {
	keys := os.Getenv("SLACK_KEYS")
	return &SlackService{
		Keys: keys,
	}
}

func (s *SlackService) AskApproval(user model.User) error {
	blockMessage := map[string]interface{}{
		"blocks": []interface{}{
			map[string]interface{}{
				"type": "header",
				"text": map[string]interface{}{
					"type": "plain_text",
					"text": "🥷 Novo usuário criado",
				},
			},
			map[string]interface{}{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("> Novo usuário criado com sucesso!\nNome: %s\nEmail: %s\n\nPor favor, aprove ou rejeite o usuário.", user.Name, user.Email),
				},
			},
			map[string]interface{}{
				"type": "actions",
				"elements": []interface{}{
					map[string]interface{}{
						"type": "button",
						"text": map[string]interface{}{
							"type": "plain_text",
							"text": "Aprovar",
						},
						"style":     "primary",
						"value":     fmt.Sprintf("approve_%s", user.ID),
						"action_id": "approve",
					},
					map[string]interface{}{
						"type": "button",
						"text": map[string]interface{}{
							"type": "plain_text",
							"text": "Rejeitar",
						},
						"style":     "danger",
						"value":     fmt.Sprintf("decline_%s", user.ID),
						"action_id": "decline",
					},
				},
			},
		},
	}

	msgBytes, err := json.Marshal(blockMessage)
	if err != nil {
		return err
	}

	request := http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   "hooks.slack.com",
			Path:   fmt.Sprintf("/services/%s", s.Keys),
		},
		Header: http.Header{
			"Content-type": []string{"application/json"},
		},
		Body: io.NopCloser(strings.NewReader(string(msgBytes))),
	}

	client := http.Client{}
	_, err = client.Do(&request)
	if err != nil {
		return err
	}

	return nil
}

func (s *SlackService) Approved(user model.User) error {
	blockMessage := map[string]interface{}{
		"blocks": []interface{}{
			map[string]interface{}{
				"type": "header",
				"text": map[string]interface{}{
					"type": "plain_text",
					"text": "✅ Usuário aprovado",
				},
			},
			map[string]interface{}{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("> O usuário %s (%s) foi aprovado. Agora ele pode acessar o sistema.", user.Name, user.Email),
				},
			},
		},
	}

	msgBytes, err := json.Marshal(blockMessage)
	if err != nil {
		return err
	}

	request := http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   "hooks.slack.com",
			Path:   fmt.Sprintf("/services/%s", s.Keys),
		},
		Header: http.Header{
			"Content-type": []string{"application/json"},
		},
		Body: io.NopCloser(strings.NewReader(string(msgBytes))),
	}

	client := http.Client{}
	_, err = client.Do(&request)
	if err != nil {
		return err
	}

	return nil
}

func (s *SlackService) Declined(user model.User) error {
	blockMessage := map[string]interface{}{
		"blocks": []interface{}{
			map[string]interface{}{
				"type": "header",
				"text": map[string]interface{}{
					"type": "plain_text",
					"text": "❌ Usuário rejeitado",
				},
			},
			map[string]interface{}{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("> O usuário %s (%s) foi rejeitado. Por favor, entre em contato com o administrador para mais informações.", user.Name, user.Email),
				},
			},
		},
	}

	msgBytes, err := json.Marshal(blockMessage)
	if err != nil {
		return err
	}

	request := http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   "hooks.slack.com",
			Path:   fmt.Sprintf("/services/%s", s.Keys),
		},
		Header: http.Header{
			"Content-type": []string{"application/json"},
		},
		Body: io.NopCloser(strings.NewReader(string(msgBytes))),
	}

	client := http.Client{}
	_, err = client.Do(&request)
	if err != nil {
		return err
	}

	return nil
}

func (s *SlackService) NotifyPullRequestCreated(user *model.User, prTitle, repoName, branchName, URL *string) error {
	blockMessage := map[string]interface{}{
		"blocks": []interface{}{
			map[string]interface{}{
				"type": "header",
				"text": map[string]interface{}{
					"type": "plain_text",
					"text": "📝 Nova Pull Request",
				},
			},
			map[string]interface{}{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("> Nova pull request aberta por %s (%s) em %s. Nesse momento passando por uma pipeline de revisão de integridade de código. Após isso será necessário a aprovação de um administrador para que seja feito o merge.", user.Name, user.Email, time.Now().Format("2006-01-02 15:04:05")),
				},
			},
			map[string]interface{}{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": "Aqui está o link para a nova Pull Request",
				},
				"accessory": map[string]interface{}{
					"type": "button",
					"text": map[string]interface{}{
						"type":  "plain_text",
						"text":  "Abrir Pull Request",
						"emoji": true,
					},
					"url": *URL,
				},
			},
		},
	}

	msgBytes, err := json.Marshal(blockMessage)
	if err != nil {
		return err
	}

	request := http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   "hooks.slack.com",
			Path:   fmt.Sprintf("/services/%s", s.Keys),
		},
		Header: http.Header{
			"Content-type": []string{"application/json"},
		},
		Body: io.NopCloser(strings.NewReader(string(msgBytes))),
	}

	client := http.Client{}
	_, err = client.Do(&request)
	if err != nil {
		return err
	}

	return nil
}
