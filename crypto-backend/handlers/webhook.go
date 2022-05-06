package handlers

import (
	"bytes"
	"context"
	"crypto-backend/db"
	"crypto-backend/models"
	"crypto-backend/utils"
	"encoding/json"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

type reqBody struct {
	CallbackUrl string `json:"callbackUrl"`
	Platform    string `json:"platform"`
}

var telegramBotConnected = false
var telegramCallbackUrl string

var messengerBotConnected = false
var messengerCallbackUrl string

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
		telegramBotConnected = true
		telegramCallbackUrl = parsedBody.CallbackUrl
	case "messenger":
		messengerBotConnected = true
		messengerCallbackUrl = parsedBody.CallbackUrl
	}

	utils.LogCreated(req)
}

func CheckForLimitPassing(coinUpdatedChan chan bool) {
	for {
		select {
		case <-coinUpdatedChan:
			coinMap := db.CreateCoinHashMap()
			checkLimits(coinMap)
		}
	}
}

func checkLimits(coinMap map[string]float64) {
	log.Println("Checking User Limits")
	now := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCursor, err := db.UsersCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Panicf("Can't create userCursor, err: %v", err)
	}

	sendMsgCounter := 0

	for userCursor.Next(ctx) {
		var user models.User
		if err = userCursor.Decode(&user); err != nil {
			log.Panicf("Can't decode user, err: %v", err)
		}

		switch user.Platform {
		case "telegram":
			if telegramBotConnected {
				for _, limit := range user.LimitList {
					if limit.IsUpper && coinMap[limit.Code] > limit.Rate ||
						!limit.IsUpper && coinMap[limit.Code] < limit.Rate {
						limit.Rate = coinMap[limit.Code]
						limitMsg := models.LimitMsg{UserId: user.Id, Limit: limit}
						go sendLimitMessage(limitMsg, &telegramBotConnected, telegramCallbackUrl)

						sendMsgCounter += 1
					}
				}
			}
			break

		case "messenger":
			if messengerBotConnected {
				for _, limit := range user.LimitList {
					if limit.IsUpper && coinMap[limit.Code] > limit.Rate {
						limit.Rate = coinMap[limit.Code]
						limitMsg := models.LimitMsg{UserId: user.Id, Limit: limit}
						go sendLimitMessage(limitMsg, &telegramBotConnected, messengerCallbackUrl)

						sendMsgCounter += 1

					}
				}
			}
			break

		default:
			log.Printf("User %s - %s doesn't have correct platform\n", user.Id, user.Name)
		}
	}

	if err = userCursor.Close(ctx); err != nil {
		log.Printf("Can't close userCursor, err: %v", err)
	}
	log.Printf("Sending %v Limit Message(s). Time took: %v", sendMsgCounter, time.Since(now))
}

func sendLimitMessage(msg models.LimitMsg, isConnected *bool, callbackUrl string) {
	sendData, err := json.Marshal(msg)
	if err != nil {
		log.Panicf("Can't marshal Limit Message, err: %v", err)
	}

	resp, err := utils.HttpClient.Post(callbackUrl+"/limits", "application/json",
		bytes.NewReader(sendData))
	if err != nil {
		log.Panicf("Can't send Limit Message, err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		*isConnected = false
		log.Printf("Incorrect response format, closing webhook to URL: %s", callbackUrl)
	}
	utils.CloseResponseBody(resp)
}
