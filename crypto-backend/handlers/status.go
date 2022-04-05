package handlers

import (
	"log"
	"net/http"
)

func StatusHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write([]byte("{}")); err != nil {
		log.Panicf("Can't send response to mongoClient, %v", err)
	}
}
