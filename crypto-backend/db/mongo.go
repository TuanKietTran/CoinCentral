package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var mongoClient *mongo.Client
var cryptoDB *mongo.Database
var err error
var mongoCtx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

var NotInitClient = errors.New("MongoClient hasn't been init, please call StartMongoClient() first")

func StartMongoClient() {
	mongoClient, err = mongo.Connect(mongoCtx, options.Client().ApplyURI("mongodb://mongodb:27017"))
	if err != nil {
		log.Fatalf("Can't connect to MongoDB, %v", err)
	}

	// Check if Client has been connected correctly
	if err := mongoClient.Ping(mongoCtx, readpref.Primary()); err != nil {
		panic(err)
	}

	cryptoDB = mongoClient.Database("crypto-db")
}

func GetMongoClient() (mongo.Client, error) {
	if mongoClient != nil {
		return *mongoClient, nil
	} else {
		return mongo.Client{}, NotInitClient
	}
}

func GetCryptoDB() (mongo.Database, error) {
	if cryptoDB != nil {
		return *cryptoDB, nil
	} else {
		return mongo.Database{}, nil
	}
}

func StopMongoClient() {
	cancel()
	if err := mongoClient.Disconnect(mongoCtx); err != nil {
		log.Fatalf("Can't close MongoDB connection, %v", err)
	}
}
