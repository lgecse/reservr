package main

import (
	"flag"
	"fmt"

	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
	log "github.com/sirupsen/logrus"
)

func main() {

	api := webexteams.NewClient()

	var tokenFlag string
	flag.StringVar(&tokenFlag, "token", "", "The token for the Webex API")

	var endpoint string
	flag.StringVar(&endpoint, "api-endpoint", "", "The The URL of the chatbot endpoint for the webhook")

	flag.Parse()

	if tokenFlag != "" {
		api.SetAuthToken(tokenFlag)
	}

	webhookList, _, err := api.Webhooks.ListWebhooks(&webexteams.ListWebhooksQueryParams{})
	if err != nil {
		log.Panic(err)
		return
	}

	found := false

	for _, item := range webhookList.Items {
		if item.Name == "reservr-webhook" && item.TargetURL == endpoint {
			found = true
		}
		if (item.Name == "reservr-webhook" && item.TargetURL != endpoint) || (item.Name != "reservr-webhook" && item.TargetURL == endpoint) {
			api.Webhooks.DeleteWebhook(item.ID)
			fmt.Printf("Webhook deleted: Name: %s -> URL: %s\n", item.Name, item.TargetURL)
		}

	}
	request := webexteams.WebhookCreateRequest{
		Name:      "reservr-webhook",
		TargetURL: endpoint,
		Resource:  "messages",
		Event:     "created",
	}

	if !found {
		createdWebhook, _, err := api.Webhooks.CreateWebhook(&request)

		fmt.Printf("Webhook created: Name: %s -> URL: %s\n", createdWebhook.Name, createdWebhook.TargetURL)
		if err != nil {
			log.Panic(err)
			return
		}
	}

}
