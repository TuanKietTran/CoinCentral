package handlers

import (
	"bytes"
	"crypto-backend/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var metaOffRequest = bytes.NewBuffer([]byte(`{
	"currency": "USD",
	"order": "ascending",
    "sort": "rank",
    "limit": 100,
    "meta": false
}`))
var metaOnRequest = bytes.NewBuffer([]byte(`{
	"currency": "USD",
	"order": "ascending",
    "sort": "rank",
    "limit": 100,
    "meta": true
}`))

var apiKey, apiExists = os.LookupEnv("APIKEY")

func FetchCrypto() {
	if !apiExists {
		log.Panicln("$APIKEY must exists")
	}

	client := &http.Client{}

	var metaList []models.CoinMeta
	metaList = fetchMeta(client)
	log.Println(metaList)

	//var valueList []models.CoinRate
	//valueList = fetchValue(client)

	rateTicker := time.NewTimer(30 * time.Second)
	defer rateTicker.Stop()

	metaTicker := time.NewTicker(5 * time.Minute)
	defer metaTicker.Stop()

	for {
		select {
		case <-rateTicker.C:

		}
	}
}

func fetchMeta(client *http.Client) []models.CoinMeta {
	req, err := http.NewRequest("POST", "https://api.livecoinwatch.com/coins/list", metaOnRequest)
	if err != nil {
		log.Panicf("Can't create new HTTP Request, %v", err)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Panicf("Can't fetch Response, %v", err)
	} else if resp.StatusCode != http.StatusOK {
		log.Panicf("Response Status Code %v", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Can't read Reponse body, %v", err)
	}

	if err = resp.Body.Close(); err != nil {
		log.Panicf("Can't close Response body, %v", err)
	}

	var metaList []models.CoinMeta
	if err = json.Unmarshal(respBody, &metaList); err != nil {
		log.Panicf("Can't parse Response body into Array of CoinMeta, %v", err)
	}

	return metaList
}

func fetchValue(client http.Client) []models.CoinRate {
	req, err := http.NewRequest("POST", "https://api.livecoinwatch.com/coins/list", metaOffRequest)
	if err != nil {
		log.Panicf("Can't create new HTTP Request, %v", err)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Panicf("Can't fetch Response, %v", err)
	} else if resp.StatusCode != http.StatusOK {
		log.Panicf("Response Status Code %v", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Can't read Reponse body, %v", err)
	}

	if err = resp.Body.Close(); err != nil {
		log.Panicf("Can't close Response body, %v", err)
	}

	var valueList []models.CoinRate
	if err = json.Unmarshal(respBody, &valueList); err != nil {
		log.Panicf("Can't parse Response body into Array of CoinRate, %v", err)
	}

	return valueList
}
