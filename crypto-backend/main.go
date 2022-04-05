package main

import (
	"crypto-backend/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	handlers.FetchCrypto()
	fmt.Println(os.Getenv("APIKEY"))

	mux := http.NewServeMux()
	mux.HandleFunc("/status", handlers.StatusHandler)

	log.Panic(http.ListenAndServe(":80", mux))
}
