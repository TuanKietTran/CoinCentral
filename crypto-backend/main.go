package main

import (
	"context"
	"crypto-backend/db"
	"crypto-backend/handlers"
	"crypto-backend/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("Can't load .env file, err: %v", err)
	}
	config := utils.ReadConfigFile("config/config.yml")

	db.StartMongoClient(config, true)
	defer db.StopMongoClient()

	// Start goroutine of fetching crypto
	go handlers.FetchCrypto(config)

	// Setup router
	router := mux.NewRouter()
	router.HandleFunc("/status", handlers.StatusHandler)

	// Coins routes
	router.HandleFunc("/coins", handlers.SupportedCoinsHandler).Methods("GET")
	router.HandleFunc("/coins/{code}", handlers.CoinCodeHandler).Methods("GET")

	// Users routes
	router.HandleFunc("/users", handlers.CreateUserHandler).Methods("POST")
	router.HandleFunc("/users/{userId}", handlers.GetUserHandler).Methods("GET")
	router.HandleFunc("/users/{userId}", handlers.DeleteUserHandler).Methods("DELETE")

	server := http.Server{
		Handler:      router,
		Addr:         ":80",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	// Setup graceful shutdown
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-termChan
		log.Println("Stopping server")
		if err := server.Shutdown(context.TODO()); err != nil {
			log.Fatalf("Can't shutdown server, %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if err.Error() != "http: Server closed" {
			log.Printf("HTTP server closed with: %v\n", err)
		}
		log.Printf("HTTP server shut down")
	}
}
