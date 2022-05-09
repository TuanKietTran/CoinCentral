package main

import (
	"fmt"
	"log"
	"net/http"
)

func Parse(w http.ResponseWriter, req *http.Request) {
	fmt.Print("====================================================")
	fmt.Print(req.Body)
	fmt.Print("===================== BODY =========================")
	fmt.Print(req.Response.Body)
}

func main() {
	http.HandleFunc("/upload", Parse)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
