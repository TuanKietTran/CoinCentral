package handlers

import (
	"context"
	"crypto-backend/db"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type supportedCoin struct {
	Code string `bson:"code" json:"code"`
	Name string `bson:"name" json:"name"`
}

func SupportedCoinsHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	returnField := bson.D{{"code", 1}, {"name", 1}}

	cursor, err := db.CoinsCollection.Find(ctx, bson.D{},
		options.Find().SetProjection(returnField))

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v - %s - %v, err:%v\n", req.Method, req.URL, http.StatusInternalServerError, err)
		return
	}

	var supportedCoinList []supportedCoin
	if err = cursor.All(ctx, &supportedCoinList); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v %s %v, err:%v\n", req.Method, req.URL, http.StatusInternalServerError, err)
		return
	}

	parsedResp, err := json.Marshal(supportedCoinList)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v %s %v, err:%v\n", req.Method, req.URL, http.StatusInternalServerError, err)
		return
	}

	if _, err = writer.Write(parsedResp); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v %s %v, err:%v\n", req.Method, req.URL, http.StatusInternalServerError, err)
		return
	}

	writer.WriteHeader(http.StatusOK)
	log.Printf("%v %s %v\n", req.Method, req.URL, http.StatusOK)
}
