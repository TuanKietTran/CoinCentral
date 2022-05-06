package db

import (
	"context"
	"crypto-backend/models"
	"crypto-backend/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client
var CryptoDB *mongo.Database

var UsersCollection *mongo.Collection
var CoinsCollection *mongo.Collection

var err error // Share error

func StartMongoClient(config *utils.Config) {
	/*
		This function start Client which connects to MongoDB
		If generateCollection == true, also generate the required collections. Default is false
	*/
	mongoUri, mongoUriExists := os.LookupEnv("MONGO_URI")
	if !mongoUriExists {
		log.Panicf("$MONGO_URI not exists")
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatalf("Can't connect to MongoDB, %v", err)
	}

	// Check if Client has been connected correctly
	if err = MongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Panicf("Can't ping MongoDB server, err: %v", err)
	}

	dbName, dbExists := os.LookupEnv("DB")
	if !dbExists {
		dbName = "crypto-dbName"
	}
	CryptoDB = MongoClient.Database(dbName)

	// Get collection pointers
	UsersCollection = CryptoDB.Collection("Users")
	CoinsCollection = CryptoDB.Collection("Coins")

	createUserIndex()

	// Check if Coins collection has been generated with enough coins
	numOfCoins, err := CoinsCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Panicf("Can't get number of Coins, err: %v", err)
	}

	if int(numOfCoins) != config.Coins.NumOfSupportingCoins {
		if err = CoinsCollection.Drop(ctx); err != nil {
			log.Panicf("Can't recreate Coins collection, err: %v", err)
		}

		initCoinsCollection(config)
	}
}

func StopMongoClient() {
	log.Println("Closing MongoDB client")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := MongoClient.Disconnect(ctx); err != nil {
		log.Fatalf("Can't close MongoDB connection, %v", err)
	}
}

func createUserIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexList, err := UsersCollection.Indexes().ListSpecifications(ctx)
	if err != nil {
		log.Panicf("Can't read User collection indexes, err: %v", err)
	}

	for _, index := range indexList {
		if index.Name == "idIndex" {
			// Index exists, stop function
			return
		}
	}

	idIndex := mongo.IndexModel{
		Keys:    bson.D{{"id", 1}, {"platform", 1}},
		Options: options.Index().SetName("idIndex"),
	}

	if _, err = UsersCollection.Indexes().CreateOne(ctx, idIndex); err != nil {
		log.Panicf("Can't create userId Index for User collection, err: %v", err)
	}
}

func initCoinsCollection(config *utils.Config) {
	/*
		This function is used to generate the collections that our project needs
	*/

	// Create simple context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate `Code` index for Coins collection
	codeIndexModel := mongo.IndexModel{Keys: bson.D{{"code", 1}},
		Options: options.Index().SetName("codeIndex").SetUnique(true)}

	if _, err := CryptoDB.Collection("Coins").Indexes().CreateOne(ctx, codeIndexModel); err != nil {
		log.Panicf("Can't create index for `code`, err: %v", err)
	}
	log.Println("codeIndex created")

	// Fill Coins collection with values
	fillCoinsCollection(config)
}

func fillCoinsCollection(config *utils.Config) {
	/*
		Init Coin Collection with coins. These coins will be our supported coins
	*/
	log.Println("Initializing Coins collection")
	now := time.Now()

	// Get $APIKEY
	apiKey, apiExists := os.LookupEnv("APIKEY")
	if !apiExists {
		log.Panicf("Live Coin Watch API key must be stored at $APIKEY")
	}

	// Request for generating Coins collection
	firstReqPayload := strings.NewReader(fmt.Sprintf(`{
	"currency": "USD",
	"order": "ascending",
	"sort": "rank",
	"limit": %v,
	"meta": true
}`, config.Coins.NumOfSupportingCoins))

	req, err := http.NewRequest("POST", "https://api.livecoinwatch.com/coins/list", firstReqPayload)
	if err != nil {
		log.Panicf("Can't create POST request, %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", apiKey)

	resp, err := utils.HttpClient.Do(req)
	defer utils.CloseResponseBody(resp)
	if err != nil {
		log.Panicf("Can't send resquest, err: %v", err)
	}

	// Parse response's body
	var coinList []models.Coin
	if err = json.NewDecoder(resp.Body).Decode(&coinList); err != nil {
		log.Panicf("Can't parse response, err: %v", err)
	}

	// Create an array of interface for `InsertMany`
	coinsInterface := make([]interface{}, len(coinList))
	for index, coin := range coinList {
		coinsInterface[index] = coin
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err = CoinsCollection.InsertMany(ctx, coinsInterface); err != nil {
		log.Panicf("Can't insert coins into collection, err: %v", err)
	}

	log.Printf("Finished init Coins collection, time took: %v", time.Since(now))
}

func CreateCoinHashMap() map[string]float64 {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := CoinsCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Panicf("Can't connect to Coin collection, err: %v", err)
	}

	coinMap := make(map[string]float64)
	var coin models.Coin
	for cursor.Next(ctx) {
		if err = cursor.Decode(&coin); err != nil {
			log.Panicf("Can't parsed result coin, err: %v", err)
		}

		coinMap[coin.Code] = coin.Rate
	}

	if err = cursor.Close(ctx); err != nil {
		log.Printf("Can't close cursor for fetching limit")
	}
	return coinMap
}
