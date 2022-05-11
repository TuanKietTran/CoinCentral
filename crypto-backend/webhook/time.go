package webhook

import (
	"bytes"
	"container/list"
	"context"
	"crypto-backend/db"
	"crypto-backend/models"
	"crypto-backend/utils"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type simplifyCoin struct {
	Code string  `json:"code"`
	Rate float64 `json:"rate"`
}

type timeMsg struct {
	UserId models.UserId  `json:"userId"`
	Coins  []simplifyCoin `json:"coins"`
}

var timeList = make(map[string]*list.List)

func TimeThread(url *URL, termChan chan bool) {
	/*
		This thread will
	*/
	var timeCounter time.Time

	for i := 0; i < 24*60; i++ {
		timeList[timeCounter.Format(time.Kitchen)] = list.New()
		timeCounter = timeCounter.Add(1 * time.Minute)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var user models.User
	cursor, err := db.UsersCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Panicf("Can't fetch Users from MongoDB, err: %v", err)
	}

	for cursor.Next(ctx) {
		if err := cursor.Decode(&user); err != nil {
			log.Panicf("Can't parse User, err: %v", err)
		}

		for _, notifyTime := range user.TimeList {
			timeList[notifyTime.Format(time.Kitchen)].PushBack(models.UserId{Id: user.Id, Platform: user.Platform})
		}
	}

	if err = cursor.Close(ctx); err != nil {
		log.Printf("Can't close cursor, just skipping, err: %v", err)
	}
	cancel()

	ticker := time.NewTicker(time.Minute)
	sendMsgCtx, cancel := context.WithCancel(context.Background())

	for {
		select {
		case <-termChan:
			log.Println("Terminating Time thread")
			cancel()
			ticker.Stop()

		case now := <-ticker.C:
			// For testing purpose: Create new item every time
			//timeList[now.Format(time.Kitchen)].PushBack(
			//	models.UserId{Id: "1972606077", Platform: "telegram"})

			go handleUserIdList(url, timeList[now.Format(time.Kitchen)], sendMsgCtx)

		case timeMsg := <-utils.TimeUpdateChan:
			switch timeMsg.Type {
			case utils.Insert:
				timeList[timeMsg.Time.Format(time.Kitchen)].PushBack(timeMsg.UserId)
			case utils.Delete:
				go popUserIdFromList(timeMsg.UserId, timeList[timeMsg.Time.Format(time.Kitchen)])
			}
		}
	}
}

func popUserIdFromList(userId models.UserId, list *list.List) {
	index := list.Front()
	for index != nil && index.Value != userId {
		index = list.Front()
	}

	if index != nil {
		list.Remove(index)
	}
}

func handleUserIdList(url *URL, userIdList *list.List, ctx context.Context) {
	now := time.Now()
	filterList := make([]interface{}, userIdList.Len())

	if userIdList.Len() == 0 {
		return
	}
  
	index := 0
	for elm := userIdList.Front(); elm != nil; elm = elm.Next() {
		userId := elm.Value.(models.UserId)
		filterList[index] = bson.M{"id": userId.Id, "platform": userId.Platform}
		index++
	}

	filter := bson.M{"$or": filterList}
	projection := bson.M{"id": 1, "platform": 1, "codeList": 1}

	mongoCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.UsersCollection.Find(mongoCtx, filter, options.Find().SetProjection(projection))
	if err != nil {
		log.Panicf("Can't fetch User collection, err: %v", err)
	}

	var user models.User
	msgCounter := 0

	for cursor.Next(mongoCtx) {
		select {
		case <-ctx.Done():
			log.Println("Stop sending Time Messages")
		default:
		}

		if err := cursor.Decode(&user); err != nil {
			log.Panicf("Can't parse user, err: %v", err)
		}

		if user.Platform == "telegram" && url.TelegramBotConnected {
			coinList := make([]simplifyCoin, len(user.CodeList))
			for i, coin := range user.CodeList {
				coinList[i] = simplifyCoin{
					Code: coin,
					Rate: db.CoinHashMap[coin],
				}
			}

			sendTimeMessage(timeMsg{
				UserId: models.UserId{
					Id:       user.Id,
					Platform: user.Platform},
				Coins: coinList}, &url.TelegramBotConnected, url.TelegramCallbackUrl)
			msgCounter++
		} else if user.Platform == "messenger" && url.MessengerBotConnected {
			coinList := make([]simplifyCoin, len(user.CodeList))
			for i, coin := range user.CodeList {
				coinList[i] = simplifyCoin{
					Code: coin,
					Rate: db.CoinHashMap[coin],
				}
			}

			sendTimeMessage(timeMsg{
				UserId: models.UserId{
					Id:       user.Id,
					Platform: user.Platform},
				Coins: coinList}, &url.MessengerBotConnected, url.MessengerCallbackUrl)
			msgCounter++
		} else {
			log.Printf("Unknown platform %v for user with ID %v", user.Platform, user.Id)
		}

	}

	if err = cursor.Close(mongoCtx); err != nil {
		log.Printf("Can't close cursor, just skipping, err: %v", err)
	}

	log.Printf("Finished sending Time messages. Sent %v Times message(s). Time took: %v",
		msgCounter, time.Since(now))
}

func sendTimeMessage(msg timeMsg, isConnected *bool, callbackUrl string) {
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
