package handlers

import (
	"net/http"
)

func StatusHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	if req.Method != http.MethodGet || req.ContentLength > 0 {
		writer.WriteHeader(http.StatusBadRequest)
	} else {
		writer.WriteHeader(http.StatusOK)
	}
}
