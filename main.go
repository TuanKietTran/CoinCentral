package main

import (
	"go-facebook-bot/pkg/fb"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/webhook", fb.HandleMessenger)
	http.HandleFunc("/", homepageHandler)

	// port := ":4000"
	log.Println(os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home page here !!"))
}
