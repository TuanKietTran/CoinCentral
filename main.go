package main

import (
	"go-facebook-bot/pkg/fb"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", fb.HandleMessenger)

	port := ":4000"
	log.Fatal(http.ListenAndServe(port, nil))
}
