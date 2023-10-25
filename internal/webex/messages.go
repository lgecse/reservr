package webex

import (
	"github.com/go-resty/resty/v2"
	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
)

type MessagesClient interface {
	GetMessage(id string) (*webexteams.Message, *resty.Response, error)
	CreateMessage(messageCreateRequest *webexteams.MessageCreateRequest) (*webexteams.Message, *resty.Response, error)
}
