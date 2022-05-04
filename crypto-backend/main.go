package main

import (
	"context"
	"crypto-backend/db"
	"crypto-backend/handlers"
	"crypto-backend/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil && err.Error() != "open .env: no such file or directory" {
		log.Panicf("Can't load .env file, err: %v", err)
	}
	config := utils.ReadConfigFile("config/config.yml")

	db.StartMongoClient(config)
	defer db.StopMongoClient()

	// Start goroutine of fetching crypto
	go handlers.FetchCrypto(config)

	// Setup router
	router := mux.NewRouter()
	router.HandleFunc("/status", handlers.StatusHandler).Methods("GET")

	// Coins routes
	router.HandleFunc("/coins", handlers.SupportedCoinsHandler).Methods("GET")
	router.HandleFunc("/coins/{code}", handlers.CoinCodeHandler).Methods("GET")

	// Users routes
	router.HandleFunc("/users", handlers.CreateUserHandler).Methods("POST")
	router.HandleFunc("/users", handlers.GetUserHandler).Methods("GET")
	router.HandleFunc("/users", handlers.DeleteUserHandler).Methods("DELETE")

	// Limit Notification routes
	router.HandleFunc("/notifications/limits", handlers.CreateLimitHandler).Methods("POST")
	router.HandleFunc("/notifications/limits", handlers.UpdateLimitHandler).Methods("PUT")
	router.HandleFunc("/notifications/limits", handlers.GetLimitHandler).Methods("GET")
	router.HandleFunc("/notifications/limits", handlers.DeleteLimitHandler).Methods("DELETE")

	// Timer Notification routes
	router.HandleFunc("/notifications/time", handlers.PostTimerHandler).Methods("POST")
	router.HandleFunc("/notifications/time", handlers.GetTimerHandler).Methods("GET")
	router.HandleFunc("/notifications/time", handlers.DeleteTimerHandler).Methods("DELETE")

	// Serve documentation
	router.
		PathPrefix("/docs").
		Handler(http.StripPrefix("/docs", http.FileServer(http.Dir("./docs/"))))

	listeningPort, listeningPortExists := os.LookupEnv("PORT")
	if !listeningPortExists {
		log.Panicf("$PORT not exists")
	}

	server := http.Server{
		Handler:      router,
		Addr:         ":" + listeningPort,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	// Setup graceful shutdown
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

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
