package main

import (
	"go-facebook-bot/pkg/fb"
	"log"
	"net/http"
	"os"
)

const projectDirName = "go-facebook-bot"

func main() {

	http.HandleFunc("/webhook", fb.HandleMessenger)
	http.HandleFunc("/", homepageHandler)

	// port := ":4000"
	log.Println(os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":8089", nil))
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home page here !!"))
}
