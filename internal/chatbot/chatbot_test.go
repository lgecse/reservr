package chatbot_test

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
	"github.com/kristofgyuracz/reservr/internal/chatbot"
	"github.com/kristofgyuracz/reservr/internal/webex"
	"github.com/stretchr/testify/assert"
)

var webexEvent1 = &webexteams.Event{ //nolint:exhaustivestruct
	Data: struct {
		ID          string "json:\"id,omitempty\""
		RoomID      string "json:\"roomId,omitempty\""
		RoomType    string "json:\"roomType,omitempty\""
		Text        string "json:\"text,omitempty\""
		PersonID    string "json:\"personId,omitempty\""
		PersonEmail string "json:\"personEmail,omitempty\""
		Created     string "json:\"created,omitempty\""
		Type        string "json:\"type,omitempty\""
	}{
		ID: "1",
	},
}

var messageCreateRequest1 = webexteams.MessageCreateRequest{
	ToPersonID: "9999@pers.on",
	Text:       "TODO#1: Help message",
}

var messageCreateRequest2 = webexteams.MessageCreateRequest{
	ToPersonID: "9999@pers.on",
	Text:       "Echoing your message back: echo cho ho",
}

var messageCreateRequest3 = webexteams.MessageCreateRequest{
	ToPersonID: "9999@pers.on",
	Text:       "Echoing your message back: ECHO cho ho o",
}

var message1 = webexteams.Message{}

func TestHandleMessage(t *testing.T) {
	t.Helper()

	tests := []struct {
		name              string
		event             *webexteams.Event
		sentMessage       *webexteams.Message
		expectedMessage   string
		expectedRecipient string
		expectedError     error
	}{
		{
			"should return with usage if command is help",
			webexEvent1,
			&webexteams.Message{
				Text:     "help",
				PersonID: "9999@pers.on",
			},
			"TODO#1: Help message",
			"9999@pers.on",
			nil,
		},
		{
			"should return with echo if command is echo",
			webexEvent1,
			&webexteams.Message{
				Text:     "echo cho ho",
				PersonID: "9999@pers.on",
			},
			"Echoing your message back: echo cho ho",
			"9999@pers.on",
			nil,
		},
		{
			"should return with usage if command is not valid",
			webexEvent1,
			&webexteams.Message{
				Text:     "invalidcommand",
				PersonID: "9999@pers.on",
			},
			"TODO#1: Help message",
			"9999@pers.on",
			nil,
		},
		{
			"should return with echo if command is ECHO",
			webexEvent1,
			&webexteams.Message{
				Text:     "ECHO cho ho o",
				PersonID: "9999@pers.on",
			},
			"Echoing your message back: ECHO cho ho o",
			"9999@pers.on",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			messagesClient := new(webex.MockMessagesClient)
			bot := chatbot.NewChatbot(messagesClient)

			messagesClient.On("GetMessage", tt.event.Data.ID).Return(tt.sentMessage, &resty.Response{}, nil)
			messagesClient.On("CreateMessage", &messageCreateRequest1).Return(&message1, &resty.Response{RawResponse: &http.Response{StatusCode: 200}}, nil)
			messagesClient.On("CreateMessage", &messageCreateRequest2).Return(&message1, &resty.Response{RawResponse: &http.Response{StatusCode: 200}}, nil)
			messagesClient.On("CreateMessage", &messageCreateRequest3).Return(&message1, &resty.Response{RawResponse: &http.Response{StatusCode: 200}}, nil)

			message, err := bot.HandleMessage(tt.event)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedMessage, message)
		})
	}
}
