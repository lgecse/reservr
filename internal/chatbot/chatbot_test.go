package chatbot_test

import (
	"fmt"
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

const ExpectedEchoMessageFormat = "Echoing your message back: %s"

const ExpectedHelpMessage = `Usage: 

Commands: 
	help  		Shows this help
	echo		Replies with your message
	reserve 	Make a reservation
	cancel		Cancel a reservation
	get			Get reservation information

Examples:
	help						Prints the help message
	echo foo bar				Replies with echo foo bar
	reserve P164				Reserves P164 parking lot for today
	reserve P164 2023-10-25		Reserves P164 parking lot for 2023-10-25
	reserve P164 tomorrow		Reserves P164 parking lot for tomorrow
	cancel T058 2023-10-26		Cancels reservation of T058 table for 2023-10-26
	get							Returns with your future reservations
	get 2023-10-27				Returns with your reservations for 2023-10-27

Resource examples:
	P164		Parking lot ID 164
	T058		Table ID 058`

var messageCreateRequest1 = webexteams.MessageCreateRequest{
	ToPersonID: "9999@pers.on",
	Text:       ExpectedHelpMessage,
}

var messageCreateRequest2 = webexteams.MessageCreateRequest{
	ToPersonID: "9999@pers.on",
	Text:       fmt.Sprintf(ExpectedEchoMessageFormat, "echo cho ho"),
}

var messageCreateRequest3 = webexteams.MessageCreateRequest{
	ToPersonID: "9999@pers.on",
	Text:       fmt.Sprintf(ExpectedEchoMessageFormat, "ECHO cho ho o"),
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
			ExpectedHelpMessage,
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
			ExpectedHelpMessage,
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
