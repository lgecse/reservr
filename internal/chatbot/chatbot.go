package chatbot

import (
	"fmt"

	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
	"github.com/kristofgyuracz/reservr/internal/command"
	"github.com/kristofgyuracz/reservr/internal/webex"
	log "github.com/sirupsen/logrus"
)

type Chatbot struct {
	messagesClient webex.MessagesClient
}

func NewChatbot(client webex.MessagesClient) *Chatbot {
	instance := Chatbot{
		messagesClient: client,
	}
	return &instance
}

func (c *Chatbot) HandleMessage(event *webexteams.Event) (string, error) {
	message, err := c.getMessageFromAPI(event)
	if err != nil {
		return "", err
	}

	var responseText string
	cmd := c.getCommand(message.Text)

	switch cmd.Command {
	case command.EchoCommand:
		responseText = c.getEchoMessage(event)
	case command.HelpCommand:
		responseText = c.getHelpMessage()
	default:
		responseText = c.getHelpMessage()
	}

	err = c.sendMessageToPerson(responseText, message.PersonID)

	return responseText, err
}

func (c *Chatbot) getHelpMessage() string {
	return `Usage: 

Commands: 
	help  		Shows this help
	echo		Replies with your message
	reserve 	Make a reservation
	cancel		Cancel a reservation
	get			Get reservation information

Examples:
	help
	echo foo bar
	reserve P164 
	reserve P164 2023-10-25
	reserve P164 tomorrow
	cancel T058 2023-10-26
	get
	get 2023-10-27

Resource examples:
	P164		Parking lot ID 164
	T058		Table ID 058`
}

func (c *Chatbot) getEchoMessage(event *webexteams.Event) string {
	message_id := event.Data.ID
	receivedMessage, _, _ := c.messagesClient.GetMessage(message_id)

	return fmt.Sprintf("Echoing your message back: %s", receivedMessage.Text)
}

func (c *Chatbot) getMessageFromAPI(event *webexteams.Event) (webexteams.Message, error) {
	message, _, err := c.messagesClient.GetMessage(event.Data.ID)
	return *message, err
}

func (c *Chatbot) getCommand(messageText string) command.CommandCall {
	commandCall := command.Parse(messageText)

	if !commandCall.IsValid() {
		commandCall = command.CommandCall{Command: command.HelpCommand}
	}

	return commandCall
}

func (c *Chatbot) sendMessageToPerson(messageText string, recipientID string) error {
	messageToSend := &webexteams.MessageCreateRequest{
		ToPersonID: recipientID,
		Text:       messageText,
	}

	_, response, err := c.messagesClient.CreateMessage(messageToSend)

	log.Debugf("Webex API responded with: %d, %s", response.StatusCode(), response.Body())
	return err
}
