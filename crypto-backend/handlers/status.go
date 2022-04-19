package handlers

import (
	"crypto-backend/utils"
	"net/http"
)

func StatusHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		writer.WriteHeader(http.StatusBadRequest)
		utils.LogBadRequest(req)
	} else {
		writer.WriteHeader(http.StatusOK)
		utils.LogSuccess(req)
	}
}
