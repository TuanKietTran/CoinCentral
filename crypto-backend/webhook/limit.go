package webhook

import (
	"bytes"
	"context"
	"crypto-backend/db"
	"crypto-backend/models"
	"crypto-backend/utils"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

type limitMsg struct {
	UserId string       `json:"userId"`
	Limit  models.Limit `json:"limit"`
}

func LimitThread(coinUpdatedChan chan bool, url *URL) {
	for {
		select {
		case <-coinUpdatedChan:
			db.CoinHashMap = db.CreateCoinHashMap()
			checkLimits(url)
		}
	}
}

func checkLimits(url *URL) {
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
			if url.TelegramBotConnected {
				for _, limit := range user.LimitList {
					if limit.IsUpper && db.CoinHashMap[limit.Code] > limit.Rate ||
						!limit.IsUpper && db.CoinHashMap[limit.Code] < limit.Rate {
						limit.Rate = db.CoinHashMap[limit.Code]
						limitMsg := limitMsg{UserId: user.Id, Limit: limit}
						go sendLimitMessage(limitMsg, &url.TelegramBotConnected, url.TelegramCallbackUrl)

						sendMsgCounter += 1
					}
				}
			}
			break

		case "messenger":
			if url.MessengerBotConnected {
				for _, limit := range user.LimitList {
					if limit.IsUpper && db.CoinHashMap[limit.Code] > limit.Rate {
						limit.Rate = db.CoinHashMap[limit.Code]
						limitMsg := limitMsg{UserId: user.Id, Limit: limit}
						go sendLimitMessage(limitMsg, &url.TelegramBotConnected, url.MessengerCallbackUrl)

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

func sendLimitMessage(msg limitMsg, isConnected *bool, callbackUrl string) {
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
