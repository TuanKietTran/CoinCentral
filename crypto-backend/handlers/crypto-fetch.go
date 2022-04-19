package handlers

import (
	"context"
	"crypto-backend/db"
	"crypto-backend/models"
	"crypto-backend/utils"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func FetchCrypto(config *utils.Config) {
	var apiKey, apiExists = os.LookupEnv("APIKEY")
	if !apiExists {
		log.Panicln("$APIKEY not exists")
	}

	fetchRankAndInsert(config, apiKey)

	ticker := time.NewTicker(time.Duration(config.Coins.TimeBetweenFetch) * time.Second)
	for {
		select {
		case <-ticker.C:
			fetchRankAndInsert(config, apiKey)
		}
	}

}

func fetchRankAndInsert(config *utils.Config, apiKey string) {
	log.Println("Updating Coins collection")

	reqPayload := strings.NewReader(fmt.Sprintf(`{
	"currency": "USD",
    "sort": "rank",
    "order": "ascending",
    "limit": %v,
    "meta": false
}`, config.Coins.NumOfFetchCoin))

	req, err := http.NewRequest("POST", "https://api.livecoinwatch.com/coins/list", reqPayload)
	if err != nil {
		log.Panicf("Can't create new Live Coin Watch request, %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", apiKey)

	resp, err := utils.HttpClient.Do(req)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Printf("Can't close response body")
		}
	}()
	if err != nil {
		log.Panicf("Can't fetch Live Coin Watch response, %v", err)
	}

	var coinList []models.Coin
	if err = json.NewDecoder(resp.Body).Decode(&coinList); err != nil {
		log.Panicf("Can't parse Live Coin Watch result, %v", err)

	}

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	counter := 0 // Count how many coins is updated
	for _, coin := range coinList {
		filter := bson.D{{"code", coin.Code}}
		update := bson.D{{"$set", bson.D{{"rate", coin.Rate}}}}
		result, err := db.CoinsCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			log.Panicf("Can't update coin: %s, %v", coin.Code, err)
		}

		if result.MatchedCount != 0 {
			counter += 1
		}
	}

	log.Printf("Finishing update Coins collection, updated %v coins\n", counter)
}
