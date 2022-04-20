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

func GetLimitHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get query parameter
	code := req.URL.Query().Get("code")
	userId := req.URL.Query().Get("userId")
	if userId == "" {
		writer.WriteHeader(http.StatusBadRequest)
		utils.LogBadRequest(req)
		return
	}

	result := db.UsersCollection.FindOne(ctx, bson.M{"_id": userId},
		options.FindOne().SetProjection(bson.M{"limitList": 1}))

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

	limitList := user.LimitList
	if code == "" {
		if err := resEncoder.Encode(limitList); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			utils.LogInternalError(req, err)
			return
		}
	} else {
		if req.URL.Query().Get("isUpper") == "" {
			writer.WriteHeader(http.StatusBadRequest)
			utils.LogBadRequest(req, errors.New("missing `isUpper`"))
			return
		}

		isUpper, err := strconv.ParseBool(req.URL.Query().Get("isUpper"))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			utils.LogBadRequest(req, errors.New("bad value of `isUpper`"))
			return
		}

		for index, limit := range limitList {
			if limit.Code == code && limit.IsUpper == isUpper {
				if err := resEncoder.Encode(limitList[index : index+1]); err != nil {
					writer.WriteHeader(http.StatusInternalServerError)
					utils.LogInternalError(req, err)
				}

				break
			}
		}
	}

	utils.LogSuccess(req)
}

func CreateLimitHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get userId
	userId := req.URL.Query().Get("userId")
	if userId == "" {
		writer.WriteHeader(http.StatusBadRequest)
		utils.LogBadRequest(req)
		return
	}

	// Parse body to get Limit
	var newLimit models.Limit
	reqDecoder := json.NewDecoder(req.Body)
	err := reqDecoder.Decode(&newLimit)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// We add new Limit and sort the array
	updateValue :=
		bson.M{
			"$push": bson.M{
				"limitList": bson.M{
					"$each": []models.Limit{newLimit},
					"$sort": bson.M{"code": 1, "isUpper": 1}},
			},
		}

	// Insert into MongoDB
	updateResult := db.UsersCollection.FindOneAndUpdate(ctx, bson.D{{"_id", userId}}, updateValue)
	if updateResult.Err() == mongo.ErrNoDocuments {
		writer.WriteHeader(http.StatusNotFound)
		utils.LogNotFound(req)
	} else if updateResult.Err() != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, updateResult.Err())
	} else {
		utils.LogSuccess(req)
	}
}

func UpdateLimitHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get userId
	userId := req.URL.Query().Get("userId")
	if userId == "" {
		writer.WriteHeader(http.StatusBadRequest)
		utils.LogBadRequest(req)
		return
	}

	// Parse body to get update Limit
	var updateLimit models.Limit
	reqDecoder := json.NewDecoder(req.Body)
	err := reqDecoder.Decode(&updateLimit)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter :=
		bson.M{"_id": userId,
			"limitList": bson.M{
				"$elemMatch": bson.M{
					"code":    updateLimit.Code,
					"isUpper": updateLimit.IsUpper,
				},
			},
		}

	// We add new Limit and sort the array
	updateValue :=
		bson.M{
			"$set": bson.M{
				"limitList.$.rate": updateLimit.Rate,
			},
		}

	// Insert into MongoDB
	updateResult := db.UsersCollection.FindOneAndUpdate(ctx, filter, updateValue)
	if updateResult.Err() == mongo.ErrNoDocuments {
		writer.WriteHeader(http.StatusNotFound)
		utils.LogNotFound(req)
	} else if updateResult.Err() != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, updateResult.Err())
	} else {
		utils.LogSuccess(req)
	}
}

func DeleteLimitHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get query parameter
	code := req.URL.Query().Get("code")
	isUpper, parseIsUpperErr := strconv.ParseBool(req.URL.Query().Get("isUpper"))
	userId := req.URL.Query().Get("userId")
	if userId == "" || code == "" || parseIsUpperErr != nil {
		writer.WriteHeader(http.StatusBadRequest)
		utils.LogBadRequest(req, errors.New("invalid or missing field"))
		return
	}

	// We add new Limit and sort the array
	removeOption :=
		bson.M{
			"$pull": bson.M{
				"limitList": bson.M{
					"code":    code,
					"isUpper": isUpper,
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