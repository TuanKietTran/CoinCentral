package main

import (
	"go-facebook-bot/pkg/fb"
	"log"
	"net/http"
	"os"
)

const projectDirName = "go-facebook-bot"

func getEnv() {
	// projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	// currentWorkDirectory, _ := os.Getwd()
	// rootPath := projectName.Find([]byte(currentWorkDirectory))

	// err := godotenv.Load(string(rootPath) + `/.env`)

	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }
}

func main() {
	// getEnv()

	http.HandleFunc("/webhook", fb.HandleMessenger)
	http.HandleFunc("/", homepageHandler)

	// port := ":4000"
	log.Printf("PORT = %v \n", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
	// log.Fatal(http.ListenAndServe(":8089", nil))
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home page here !!"))
}
