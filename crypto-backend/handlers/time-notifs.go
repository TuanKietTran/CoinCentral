package handlers

import (
	"context"
	"crypto-backend/db"
	"crypto-backend/models"
	"crypto-backend/utils"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"time"
)

type timeResponse struct {
	CodeList []string `json:"codeList"`
	TimeList []string `json:"timeList"`
}

func GetTimerHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId, err := utils.GetUserId(req)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		utils.LogBadRequest(req, err)
		return
	}

	// Get query
	var getCode bool
	if req.URL.Query().Get("getCode") == "" {
		getCode = false
	} else {
		getCode, err = strconv.ParseBool(req.URL.Query().Get("getCode"))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			utils.LogBadRequest(req, err)
			return
		}
	}

	var getTime bool
	if req.URL.Query().Get("getTime") == "" {
		getTime = false
	} else {
		getTime, err = strconv.ParseBool(req.URL.Query().Get("getTime"))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			utils.LogBadRequest(req, err)
			return
		}
	}

	result := db.UsersCollection.FindOne(ctx, bson.M{"id": userId.Id, "platform": userId.Platform},
		options.FindOne().SetProjection(bson.M{"codeList": 1, "timeList": 1}))

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

	var responseObj timeResponse
	if getCode {
		responseObj.CodeList = user.CodeList
	} else {
		responseObj.CodeList = []string{}
	}

	if getTime {
		responseObj.TimeList = make([]string, len(user.TimeList))
		for i, rawTime := range user.TimeList {
			responseObj.TimeList[i] = rawTime.Format(time.Kitchen)
		}
	} else {
		responseObj.TimeList = []string{}
	}

	// Encode response
	if err := json.NewEncoder(writer).Encode(responseObj); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	utils.LogSuccess(req)
}

func PostTimerHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	userId, err := utils.GetUserId(req)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		utils.LogBadRequest(req, err)
		return
	}

	updateValue := bson.M{
		"$addToSet": bson.M{},
	}

	insertFields := bson.M{}

	newTime := req.URL.Query().Get("time")
	if newTime != "" {
		parseTime, err := time.Parse(time.Kitchen, newTime)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			utils.LogBadRequest(req, err)
			return
		}

		insertFields["timeList"] = parseTime
	}

	newCode := req.URL.Query().Get("code")
	if newCode != "" {
		insertFields["codeList"] = newCode
	}

	updateValue["$addToSet"] = insertFields

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateResult := db.UsersCollection.FindOneAndUpdate(ctx,
		bson.M{"id": userId.Id, "platform": userId.Platform}, updateValue)
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

func DeleteTimerHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	userId, err := utils.GetUserId(req)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		utils.LogBadRequest(req, err)
		return
	}

	updateValue := bson.M{
		"$pull": bson.M{},
	}

	insertFields := bson.M{}

	newTime := req.URL.Query().Get("time")
	if newTime != "" {
		parseTime, err := time.Parse(time.Kitchen, newTime)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			utils.LogBadRequest(req, err)
			return
		}

		insertFields["timeList"] = parseTime
	}

	newCode := req.URL.Query().Get("code")
	if newCode != "" {
		insertFields["codeList"] = newCode
	}

	updateValue["$pull"] = insertFields

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateResult := db.UsersCollection.FindOneAndUpdate(ctx,
		bson.M{"id": userId.Id, "platform": userId.Platform}, updateValue)

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
