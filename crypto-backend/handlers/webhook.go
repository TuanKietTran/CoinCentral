package handlers

import (
	"crypto-backend/utils"
	"crypto-backend/webhook"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

type reqBody struct {
	CallbackUrl string `json:"callbackUrl"`
	Platform    string `json:"platform"`
}

var GlobalURL *webhook.URL

func CreateWebhookHandler(writer http.ResponseWriter, req *http.Request) {
	var parsedBody reqBody
	if err := json.NewDecoder(req.Body).Decode(&parsedBody); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		utils.LogBadRequest(req, err)
		return
	}

	// Try pinging the callbackUrl
	now := time.Now()
	if resp, err := utils.HttpClient.Get(parsedBody.CallbackUrl + "/ping"); err != nil || resp.StatusCode != 200 {
		utils.CloseResponseBody(resp)
		writer.WriteHeader(http.StatusBadRequest)
		if err != nil {
			utils.LogBadRequest(req, err)
		} else {
			utils.LogBadRequest(req, errors.New("response code is not 200"))
		}
		return
	}
	log.Printf("Ping time: %v", time.Since(now))

	switch parsedBody.Platform {
	case "telegram":
		GlobalURL.TelegramBotConnected = true
		GlobalURL.TelegramCallbackUrl = parsedBody.CallbackUrl
	case "messenger":
		GlobalURL.MessengerBotConnected = true
		GlobalURL.MessengerCallbackUrl = parsedBody.CallbackUrl
	}

	utils.LogCreated(req)
}
