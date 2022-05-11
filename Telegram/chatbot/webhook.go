package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/TuanKietTran/CoinCentral/object"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	print("====LIMIT HANDLER CALL====")
	w.WriteHeader(http.StatusOK)
}

func LimitsHandler(w http.ResponseWriter, r *http.Request) {
	print("====LIMIT HANDLER CALL====")
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(reqBody))
	var limitMsg object.WebhookLimitMsg
	err := json.Unmarshal(reqBody, &limitMsg)
	if err != nil {
		panic(err)
	}
	//TODO: Handle data
	limitMsgPostBacker(limitMsg)
	w.WriteHeader(http.StatusOK)
}

func TimesHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var timeMsg object.WebhookTimeMsg
	json.Unmarshal(reqBody, &timeMsg)
	//TODO: Handle data
	timeMsgPostBacker(timeMsg)
	w.WriteHeader(http.StatusOK)
}
