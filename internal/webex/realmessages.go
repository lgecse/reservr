package webex

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
)

type realMessagesClient struct {
	client *webexteams.Client
}

func New(client *webexteams.Client) MessagesClient {
	instance := realMessagesClient{
		client: client,
	}
	return &instance
}

func (m *realMessagesClient) GetMessage(id string) (*webexteams.Message, *resty.Response, error) {
	message, response, err := m.client.Messages.GetMessage(id)
	if err != nil {
		return message, response, fmt.Errorf("failed to get message: %w", err)
	}
	return message, response, nil
}

func (m *realMessagesClient) CreateMessage(messageCreateRequest *webexteams.MessageCreateRequest) (*webexteams.Message, *resty.Response, error) {
	message, response, err := m.client.Messages.CreateMessage(messageCreateRequest)
	if err != nil {
		return message, response, fmt.Errorf("failed to create message: %w", err)
	}
	return message, response, nil
}
