package utils

import (
	"log"
	"net/http"
)

func LogSuccess(req *http.Request) {
	log.Printf("%v %s 200\n", req.Method, req.URL)
}

func LogCreated(req *http.Request) {
	log.Printf("%v %s 201\n", req.Method, req.URL)
}

func LogBadRequest(req *http.Request, err ...error) {
	log.Printf("%v %s 400, err: %v", req.Method, req.URL, err)
}

func LogNotFound(req *http.Request) {
	log.Printf("%v %s 404\n", req.Method, req.URL)
}

func LogInternalError(req *http.Request, err error) {
	log.Printf("%v %s 500, err: %v\n", req.Method, req.URL, err)
}
