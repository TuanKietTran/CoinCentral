package main

import (
	"context"
	"crypto-backend/db"
	"crypto-backend/handlers"
	"crypto-backend/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config := utils.ReadConfigFile("config/config.yml")

	db.StartMongoClient(config, false)
	defer db.StopMongoClient()

	// Start goroutine of fetching crypto
	go handlers.FetchCrypto(config)

	// Setup router
	router := mux.NewRouter()
	router.HandleFunc("/status", handlers.StatusHandler)

	// Users routes
	router.HandleFunc("/users", handlers.GetUserHandler).Methods("GET")
	router.HandleFunc("/users", handlers.CreateUserHandler).Methods("POST")
	router.HandleFunc("/users", handlers.DeleteUserHandler).Methods("DELETE")

	// Coins routes
	router.HandleFunc("/coins/supported", handlers.SupportedCoinsHandler).Methods("GET")

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
