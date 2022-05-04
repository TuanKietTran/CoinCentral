package handlers

import (
	"context"
	"crypto-backend/db"
	"crypto-backend/models"
	"crypto-backend/utils"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"time"
)

func GetTimerHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := req.URL.Query().Get("userId")
	if userId == "" {
		writer.WriteHeader(http.StatusBadRequest)
		utils.LogBadRequest(req)
		return
	}

	result := db.UsersCollection.FindOne(ctx, bson.M{"_id": userId},
		options.FindOne().SetProjection(bson.M{"watchList": 1}))

	var user models.User

	if result.Err() == mongo.ErrNoDocuments {
		writer.WriteHeader(http.StatusNotFound)
		utils.LogNotFound(req)
		return
	}

	if err := result.Decode(&user); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	// Create JSON encoder
	resEncoder := json.NewEncoder(writer)
	watchList := user.WatchList

	if err := resEncoder.Encode(watchList); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	utils.LogSuccess(req)
}

func PostTimerHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := req.URL.Query().Get("userId")
	if userId == "" {
		writer.WriteHeader(http.StatusBadRequest)
		utils.LogBadRequest(req)
		return
	}

	var newNotification models.Notification
	reqDecoder := json.NewDecoder(req.Body)
	err := reqDecoder.Decode(&newNotification)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	filter :=
		bson.M{"_id": userId,
			"watchList": bson.M{
				"$elemMatch": bson.M{
					"code": newNotification.Code,
					"time": newNotification.Time,
				},
			},
		}

	newVal :=
	bson.M{
		"$set": bson.M{
			"limitList.$.rate": newNotification.Time,
		},
	}
	response := db.UsersCollection.FindOneAndUpdate(ctx, filter, newVal)

	if response.Err() == mongo.ErrNoDocuments {
		writer.WriteHeader(http.StatusNotFound)
		utils.LogNotFound(req)
	} else if response.Err() != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, response.Err())
	} else {
		utils.LogSuccess(req)
	}
}

func DeleteTimerHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get query parameter
	userId := req.URL.Query().Get("userId")

	// We add new Limit and sort the array
	removeOption :=
		bson.M{
			"$pull": bson.M{
				"watchList.$.time": bson.M{
					"time": ,
				},
			},
		}

	result := db.UsersCollection.FindOneAndUpdate(ctx, bson.M{"_id": userId}, removeOption)

	if result.Err() == mongo.ErrNoDocuments {
		writer.WriteHeader(http.StatusNotFound)
		utils.LogNotFound(req)
	} else if result.Err() != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, result.Err())
	} else {
		writer.WriteHeader(http.StatusOK)
		utils.LogSuccess(req)
	}
}