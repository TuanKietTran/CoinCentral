package handlers

import (
	"context"
	"crypto-backend/db"
	"crypto-backend/models"
	"crypto-backend/utils"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetUserHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	collection := db.UsersCollection

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get userId & platform
	userId, err := utils.GetUserId(req)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		utils.LogBadRequest(req, err)
		return
	}

	// Search for user with userId
	cursor, err := collection.Find(ctx, bson.M{"id": userId.Id, "platform": userId.Platform})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	var resultList []models.User
	if err := cursor.All(ctx, &resultList); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	if len(resultList) > 1 {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	} else if len(resultList) == 0 {
		writer.WriteHeader(http.StatusNotFound)
		utils.LogNotFound(req)
		return
	}

	// Parsing result
	respBody, err := json.Marshal(resultList[0])
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	if _, err = writer.Write(respBody); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	utils.LogSuccess(req)
}

func CreateUserHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	collection := db.UsersCollection

	// Parse body to get User info
	var newUser models.User
	err := json.NewDecoder(req.Body).Decode(&newUser)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	// Check if ID is empty
	if newUser.Id == "" {
		writer.WriteHeader(http.StatusBadRequest)
		utils.LogBadRequest(req)
		return
	}

	// Check if platform is correct
	if newUser.Platform != "telegram" && newUser.Platform != "messenger" {
		writer.WriteHeader(http.StatusBadRequest)
		utils.LogBadRequest(req)
		return
	}

	// Fields that we must also include
	newUser.LimitList = []models.Limit{}
	newUser.CodeList = []string{}
	newUser.TimeList = []time.Time{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert into MongoDB
	_, err = collection.InsertOne(ctx, newUser)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
	} else {
		writer.WriteHeader(http.StatusCreated)
		utils.LogCreated(req)
	}
}

func DeleteUserHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	collection := db.UsersCollection

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get userId & platform
	userId, err := utils.GetUserId(req)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		utils.LogBadRequest(req, err)
		return
	}

	// Search for user with userId
	deleteResult, err := collection.DeleteOne(ctx, bson.M{"id": userId.Id, "platform": userId.Platform})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	if deleteResult.DeletedCount == 0 {
		writer.WriteHeader(http.StatusNotFound)
		utils.LogNotFound(req)
	} else {
		writer.WriteHeader(http.StatusOK)
		utils.LogSuccess(req)
	}
}
