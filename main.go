package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
	"github.com/joho/godotenv"
	"github.com/kristofgyuracz/reservr/internal/chatbot"
	"github.com/kristofgyuracz/reservr/internal/webex"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fastjson"
)

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	ApiResponse := events.APIGatewayProxyResponse{}
	switch request.HTTPMethod {
	case "GET":
		ApiResponse = events.APIGatewayProxyResponse{Body: "Nothing to see here...", StatusCode: 200}
		log.Info("GET endpoint called")
	case "POST":
		err := fastjson.Validate(request.Body)
		if err != nil {
			body := "Error: Invalid JSON payload ||| " + fmt.Sprint(err) + " Body Obtained" + "||||" + request.Body
			log.Error(body)
			ApiResponse = events.APIGatewayProxyResponse{Body: body, StatusCode: 500}
			return ApiResponse, nil
		}

		webexEvent := webexteams.Event{}
		err = json.Unmarshal([]byte(request.Body), &webexEvent)
		if err != nil {
			body := "Error: Failed to unmarshal JSON request data to webex event ||| " + fmt.Sprint(err) + " Body Obtained" + "||||" + request.Body
			log.Error(body)
			ApiResponse = events.APIGatewayProxyResponse{Body: body, StatusCode: 500}
			return ApiResponse, nil
		}

		if webexEvent.Event == "created" && webexEvent.Resource == "messages" {
			chatbot := chatbot.NewChatbot(webex.New(webexteams.NewClient()))
			_, err := chatbot.HandleMessage(&webexEvent)

			if err != nil {
				body := "Error: Failed to produce response message ||| " + fmt.Sprint(err) + " Body Obtained" + "||||" + request.Body
				log.Error(body)
				ApiResponse = events.APIGatewayProxyResponse{Body: body, StatusCode: 500}
				return ApiResponse, nil
			}

			body := "Success: Message sent."
			log.Info(body)
			ApiResponse = events.APIGatewayProxyResponse{Body: body, StatusCode: 200}
		}

	}
	return ApiResponse, nil
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	e := godotenv.Load()
	if e != nil {
		log.Error(e)
	}

	lambda.Start(HandleRequest)
}
