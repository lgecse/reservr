package webex

import (
	"github.com/go-resty/resty/v2"
	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
	"github.com/stretchr/testify/mock"
)

type MockMessagesClient struct {
	mock.Mock
}

func (m *MockMessagesClient) GetMessage(id string) (*webexteams.Message, *resty.Response, error) {
	args := m.Called(id)
	return args.Get(0).(*webexteams.Message), args.Get(1).(*resty.Response), args.Error(2) //nolint:wrapcheck
}

func (m *MockMessagesClient) CreateMessage(messageCreateRequest *webexteams.MessageCreateRequest) (*webexteams.Message, *resty.Response, error) {
	args := m.Called(messageCreateRequest)
	return args.Get(0).(*webexteams.Message), args.Get(1).(*resty.Response), args.Error(2) //nolint:wrapcheck
}
