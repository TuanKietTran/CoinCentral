package utils

import (
	"crypto-backend/models"
	"errors"
	"net/http"
)

func GetUserId(req *http.Request) (models.UserId, error) {
	userId := req.URL.Query().Get("id")
	if userId == "" {
		return models.UserId{}, errors.New("id field missing")
	}

	platform := req.URL.Query().Get("platform")
	if platform == "" {
		return models.UserId{}, errors.New("platform field missing")
	}

	return models.UserId{Id: userId, Platform: platform}, nil
}
