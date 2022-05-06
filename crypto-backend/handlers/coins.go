package handlers

import (
	"context"
	"crypto-backend/db"
	"crypto-backend/models"
	"crypto-backend/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

func SupportedCoinsHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	returnField := bson.D{{"code", 1}}

	cursor, err := db.CoinsCollection.Find(ctx, bson.D{},
		options.Find().SetProjection(returnField))

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v - %s - %v, err:%v\n", req.Method, req.URL, http.StatusInternalServerError, err)
		return
	}

	var supportedCoinList []models.Coin
	if err = cursor.All(ctx, &supportedCoinList); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v %s %v, err:%v\n", req.Method, req.URL, http.StatusInternalServerError, err)
		return
	}

	resultList := make([]string, len(supportedCoinList))
	for i, coin := range supportedCoinList {
		resultList[i] = coin.Code
	}

	parsedResp, err := json.Marshal(resultList)
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

	log.Printf("%v %s %v\n", req.Method, req.URL, http.StatusOK)
}

func CoinCodeHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resultQuery := db.CoinsCollection.FindOne(ctx, bson.D{{"code", mux.Vars(req)["code"]}})
	if resultQuery.Err() == mongo.ErrNoDocuments {
		writer.WriteHeader(http.StatusNotFound)
		utils.LogNotFound(req)
		return
	}

	var resultCoin models.Coin
	err := resultQuery.Decode(&resultCoin)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	writeOutput, err := json.Marshal(resultCoin)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	_, err = writer.Write(writeOutput)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	utils.LogSuccess(req)
}
