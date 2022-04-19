package db

import (
	"context"
	"crypto-backend/models"
	"crypto-backend/utils"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var MongoClient *mongo.Client
var CryptoDB *mongo.Database

var UsersCollection *mongo.Collection
var CoinsCollection *mongo.Collection

var err error // Share error
var mongoCtx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

func StartMongoClient(config *utils.Config, generateCollection bool) {
	/*
		This function start Client connected to MongoDB
		If generateCollection == true, also generate the required collections. Default is false
	*/
	mongoUri, mongoUriExists := os.LookupEnv("MONGO_URI")
	if !mongoUriExists {
		log.Panicf("$MONGO_URI not exists")
	}
	MongoClient, err = mongo.Connect(mongoCtx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatalf("Can't connect to MongoDB, %v", err)
	}

	// Check if Client has been connected correctly
	if err = MongoClient.Ping(mongoCtx, readpref.Primary()); err != nil {
		panic(err)
	}

	CryptoDB = MongoClient.Database("crypto-db")

	// Get collection pointers
	if generateCollection == true {
		GenerateCollections(config)
	}

	UsersCollection = CryptoDB.Collection("Users")
	CoinsCollection = CryptoDB.Collection("Coins")
}

func StopMongoClient() {
	log.Println("Closing MongoDB client")
	cancel()
	if err := MongoClient.Disconnect(mongoCtx); err != nil {
		log.Fatalf("Can't close MongoDB connection, %v", err)
	}
}

func GenerateCollections(config *utils.Config) {
	/*
		This function is used to generate the collections that our project needs
	*/
	listOfCollection, err := CryptoDB.ListCollectionNames(mongoCtx, bson.D{})
	if err != nil {
		log.Panicf("Can't get list of collection names, %v", err)
	}

	requireCollections := []string{"Users", "Coins"}

	for _, collection := range requireCollections {
		if !contain(listOfCollection, collection) {
			if err = CryptoDB.CreateCollection(mongoCtx, collection); err != nil {
				log.Panicf("Can't create collection '%s', %v", collection, err)
			}
		}
	}

	// Create simple context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate `Code` index for Coins collection
	codeIndexModel := mongo.IndexModel{Keys: bson.D{{"code", 1}},
		Options: options.Index().SetName("codeIndex").SetUnique(true)}

	coinCollectionIndex := CryptoDB.Collection("Coins").Indexes()
	indexName, err := coinCollectionIndex.CreateOne(ctx, codeIndexModel)
	if err != nil {
		log.Panicf("Can't create index for `code`, %v", err)
	}
	log.Printf("%s created", indexName)

	// Fill Coins collection with values
	if !contain(listOfCollection, "Coins") {
		initCoinsCollection(config)
	}
}

func contain(stringList []string, item string) bool {
	// Check if a list contains a string
	for _, val := range stringList {
		if val == item {
			return true
		}
	}
	return false
}

func initCoinsCollection(config *utils.Config) {
	/*
		Init Coin Collection with coins. These coins will be our supported coins
	*/
	log.Println("Initializing Coins collection")

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
	if err != nil {
		log.Panicf("Can't send resquest, %v", err)
	}

	firstRespBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Can't read response body, %v", err)
	}

	// Closing response connection
	if err = resp.Body.Close(); err != nil {
		log.Printf("Can't close response connection, just skipping, %v\n", err)
	}

	var listOfCoins []models.Coin
	if err = json.Unmarshal(firstRespBody, &listOfCoins); err != nil {
		log.Panicf("Can't parse response, %v", err)
	}

	for _, coin := range listOfCoins {
		if _, err = CryptoDB.Collection("Coins").InsertOne(mongoCtx, coin); err != nil {
			log.Panicf("Can't insert coin %s into collection, %v", coin.Code, err)
		}
	}

	log.Println("Finished init Coins collection")
}
