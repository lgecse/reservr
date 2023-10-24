package chatbot

import (
	"fmt"
	"strings"

	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
	log "github.com/sirupsen/logrus"
)

type Chatbot struct {
	client *webexteams.Client
}

func GetChatbot() *Chatbot {
	instance := new(Chatbot)
	instance.client = webexteams.NewClient()
	return instance
}

func (c *Chatbot) HandleMessage(event *webexteams.Event) error {
	message, err := c.getMessageFromAPI(event)

	if err != nil {
		return err
	}

	var responseText string

	switch command := c.getCommand(message.Text); command {
	case "help":
		responseText = c.getHelpMessage()

	case "echo":
		responseText = c.getEchoMessage(event)

	default:
		responseText = c.getHelpMessage()
	}
	err = c.sendMessageToPerson(responseText, message.PersonID)

	return err
}

func (c *Chatbot) getHelpMessage() string {
	return "TODO#1: Help message"
}

func (c *Chatbot) getEchoMessage(event *webexteams.Event) string {
	message_id := event.Data.ID
	receivedMessage, _, _ := c.client.Messages.GetMessage(message_id)

	return fmt.Sprintf("Echoing your message back: %s", receivedMessage.Text)
}

func (c *Chatbot) getMessageFromAPI(event *webexteams.Event) (webexteams.Message, error) {
	message, _, err := c.client.Messages.GetMessage(event.Data.ID)
	return *message, err
}

func (c *Chatbot) getCommand(messageText string) string {
	words := strings.Split(messageText, " ")
	if len(words) > 0 {
		return words[0]
	} else {
		return "help"
	}
}

func (c *Chatbot) sendMessageToPerson(messageText string, recipientID string) error {

	messageToSend := &webexteams.MessageCreateRequest{
		ToPersonID: recipientID,
		Text:       messageText,
	}

	_, response, err := c.client.Messages.CreateMessage(messageToSend)

	log.Debugf("Webex API responded with: %d, %s", response.StatusCode(), response.Body())
	return err
}
