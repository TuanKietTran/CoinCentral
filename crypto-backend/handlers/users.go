package handlers

import (
	"context"
	"crypto-backend/db"
	"crypto-backend/models"
	"crypto-backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

func GetUserHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	collection := db.UsersCollection

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get userId from path
	userId, ok := mux.Vars(req)["userId"]
	if !ok {
		utils.LogBadRequest(req, errors.New("missing userId"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Search for user with userId
	cursor, err := collection.Find(ctx, bson.D{{"_id", userId}})
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
	reqDecoder := json.NewDecoder(req.Body)
	err := reqDecoder.Decode(&newUser)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
		return
	}

	newUser.UserId = uuid.NewString()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert into MongoDB
	result, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
	} else if _, err = writer.Write([]byte(fmt.Sprintf(`{"userId": "%s"}`, result.InsertedID))); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		utils.LogInternalError(req, err)
	} else {
		utils.LogSuccess(req)
	}
}

func DeleteUserHandler(writer http.ResponseWriter, req *http.Request) {
	log.Println("Deleting")
	writer.Header().Set("Content-Type", "application/json")
	collection := db.UsersCollection

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get userId from path
	userId, ok := mux.Vars(req)["userId"]
	if !ok {
		utils.LogBadRequest(req, errors.New("missing userId"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Search for user with userId
	deleteResult, err := collection.DeleteOne(ctx, bson.D{{"_id", userId}})
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
