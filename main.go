package main

import (
	"go-facebook-bot/pkg/fb"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", fb.HandleMessenger)

	port := ":8099"
	log.Fatal(http.ListenAndServe(port, nil))
}
