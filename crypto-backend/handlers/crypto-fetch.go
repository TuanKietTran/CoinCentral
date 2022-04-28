package handlers

import (
	"context"
	"crypto-backend/db"
	"crypto-backend/models"
	"crypto-backend/utils"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// LatestCoins List of the latest rates
var LatestCoins models.CoinList

func FetchCrypto(config *utils.Config) {
	var apiKey, apiExists = os.LookupEnv("APIKEY")
	if !apiExists {
		log.Panicln("$APIKEY not exists")
	}

	// Get list of supported Coins
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := db.CoinsCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Panicf("Can't fetch supported Coins from MongoDB, err: %v", err)
	}
	cancel()

	if err = cursor.All(ctx, &LatestCoins); err != nil {
		log.Panicf("Can't parsed list of supported Coins")
	}

	// Start operation loop
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
	defer utils.CloseResponseBody(resp)
	if err != nil {
		log.Panicf("Can't fetch Live Coin Watch response, %v", err)
	}

	// Parsed response's body
	var updatedCoins models.CoinList
	if err = json.NewDecoder(resp.Body).Decode(&updatedCoins); err != nil {
		log.Panicf("Can't parse Live Coin Watch result, %v", err)
	}
	sort.Sort(updatedCoins)

	// Create a writeModel
	writeModel := make([]mongo.WriteModel, config.Coins.NumOfSupportingCoins)
	counter := 0 // Count how many coins is updated
	updatedIndex := 0
	for i, value := range LatestCoins {
		for updatedIndex < len(updatedCoins) && updatedCoins[updatedIndex].Code < value.Code {
			updatedIndex += 1
		}

		// Max index reached, stop
		if updatedIndex == len(updatedCoins) {
			break
		}

		if value.Code == updatedCoins[updatedIndex].Code {
			writeModel[i] = mongo.NewUpdateOneModel().
				SetFilter(bson.M{"code": value.Code}).
				SetUpdate(bson.M{"$set": bson.M{
					"rate": updatedCoins[updatedIndex].Rate,
				}})
			counter += 1
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := db.CoinsCollection.BulkWrite(ctx, writeModel, options.BulkWrite().SetOrdered(false)); err != nil {
		log.Panicf("Can't update Coins collection, err: %v", err)
	}

	log.Printf("Finishing update Coins collection, updated %v coins\n", counter)
}
