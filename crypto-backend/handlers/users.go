package handlers

import (
	"context"
	"crypto-backend/db"
	"crypto-backend/models"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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

	// Check if query is correct
	query := req.URL.Query()
	if !query.Has("userId") {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Search for user with userId
	cursor, err := collection.Find(ctx, bson.D{{"_id", query.Get("userId")}})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicf("Can't get user info from DB, %v", err)
	}

	var resultList []models.User
	if err := cursor.All(ctx, &resultList); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicf("Can't parsed cursor into list of users, %v", err)
	}

	if len(resultList) > 1 {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicf("Result should not have more than 1 record, %v", err)
	} else if len(resultList) == 0 {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// Parsing result
	respBody, err := json.Marshal(resultList[0])
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicf("Can't parse object into byte, %v", err)
	}

	writer.WriteHeader(http.StatusOK)
	if _, err = writer.Write(respBody); err != nil {
		log.Panicf("Can't send response, %v", err)
	}
}

func CreateUserHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	collection := db.UsersCollection

	// Create a new User
	query := req.URL.Query()
	if !query.Has("name") {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	newUser := models.User{UserId: uuid.NewString(), Name: query.Get("name")}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert into MongoDB
	result, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		writer.WriteHeader(http.StatusOK)
		if _, err = writer.Write([]byte(fmt.Sprintf(`{"userId": "%s"}`, result.InsertedID))); err != nil {
			log.Panicf("Can't send response, %v", err)
		}
	}
}

func DeleteUserHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	collection := db.UsersCollection

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if query is correct
	query := req.URL.Query()
	if !query.Has("userId") {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Search for user with userId
	deleteResult, err := collection.DeleteOne(ctx, bson.D{{"_id", query.Get("userId")}})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicf("Can't delete user from DB, %v", err)
	}

	if deleteResult.DeletedCount == 0 {
		writer.WriteHeader(http.StatusNotFound)
	} else {
		writer.WriteHeader(http.StatusOK)
	}
}
