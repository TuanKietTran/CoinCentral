package fb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	// "io"
	// "github.com/joho/godotenv"
	"golang-messenger-chatbot/api"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Facebook credentials. It's better to store it in your secret storage.

// var (
// 	verifyToken = "random-verify-token"
// 	appSecret   = "f74c28519b1c720e849c051f6bd8b903"
// 	accessToken = "EAAPSXID9D9gBAGfFZBAWKxZB5vUQLZAZBWiFpe3nJLWrbcGiEt1HcN5FYAGNPao5z9sfIlPVgDV8oFAxNsBRIBgFZBmU8SZBan5mKfX81lVdTS1aFSsCD7GzFX04ZAFO5U899g9YYpikZCwZAyw6SSZAofkJa3KrherJHUHeJNF4cKXdT23cu4gEMov3Wwpb2HL9qr4zQt1SeYpQZDZD"
// )

var (
	verifyToken = os.Getenv("VERIFY_TOKEN")
	appSecret   = os.Getenv("APP_SECRET")
	accessToken = os.Getenv("ACCESS_TOKEN")
)

// errors
var (
	errUnknownWebHookObject = errors.New("unknown web hook object")
	errNoMessageEntry       = errors.New("there is no message entry")
)

// var coinListTemp = []string{"Ethereum", "Bitcoin", "Spy"}

var userRequest = make(map[string][]string)
var followedCoins = make(map[string][]string)
var allCoins = api.GetAllCoins()

// var followedCoins, _ = api.GetFollowCoins(recipientID, MESSENGER)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// HandleMessenger handles all incoming webhooks from Facebook Messenger.
func HandleMessenger(w http.ResponseWriter, r *http.Request) {
	log.Println("---- VERIFY TOKEN = ", verifyToken)
	if r.Method == http.MethodGet {
		HandleVerification(w, r)
		return
	}

	HandleWebHook(w, r)
}

// HandleVerification handles the verification request from Facebook.
func HandleVerification(w http.ResponseWriter, r *http.Request) {
	// log.Printf("Verify Token %v", r)
	q := r.URL.Query()
	if verifyToken != q.Get("hub.verify_token") {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.URL.Query().Get("hub.challenge")))
}

// HandleWebHook handles a webhook incoming from Facebook.
func HandleWebHook(w http.ResponseWriter, r *http.Request) {
	// log.Printf("VERIFY Token = %s \n", verifyToken)
	// log.Printf("Access Token = %s \n", accessToken)
	err := Authorize(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		log.Println("authorize", err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		log.Println("read webhook body", err)
		return
	}

	wr := WebHookRequest{}
	err = json.Unmarshal(body, &wr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		log.Println("unmarshal request", err)
		return
	}

	err = handleWebHookRequest(wr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal"))
		log.Println("handle webhook request", err)
		return
	}

	// Facebook waits for the constant message to get that everything is OK
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("EVENT_RECEIVED"))
}

func handleWebHookRequest(r WebHookRequest) error {
	if r.Object != "page" {
		return errUnknownWebHookObject
	}

	for _, we := range r.Entry {
		err := handleWebHookRequestEntry(we)
		if err != nil {
			return fmt.Errorf("handle webhook request entry: %w", err)
		}
	}

	return nil
}

func handleWebHookRequestEntry(we WebHookRequestEntry) error {
	if len(we.Messaging) == 0 { // Facebook claims that the arr always contains a single item but we don't trust them :)
		return errNoMessageEntry
	}

	em := we.Messaging[len(we.Messaging)-1]
	// log.Println("-------------Messaing = ", em)
	// log.Println("-------------Senderid = ", em.Sender.ID)
	// log.Println("---------Recipient = ", em.Recipient.ID)

	if em.Postback != nil {
		log.Println("<POSTBACK  +++   POSTBACK>")
		HandlePostBack(em.Sender.ID, em.Postback.Payload)
	}
	if em.Message != nil {
		// err := handleMessage(em.Sender.ID, em.Message.Text)
		// if err != nil {
		// 	return fmt.Errorf("handle message: %w", err)
		// }
		if em.Sender.ID != "111411174878132" {
			err := handleMessage(em.Sender.ID, em.Message.Text)
			if err != nil {
				return fmt.Errorf("handle message: %w", err)
			}
		}
	}

	return nil
}

func handleMessage(recipientID, msgText string) error {
	msgText = strings.TrimSpace(msgText)
	fmt.Printf("handle Message step has recipientID = %s \n", recipientID)

	var responseText string

	log.Printf("user text = %v", msgText)
	if recipientID == "111411174878132" {
		log.Println("<><><><><><><>")
	}

	if msgText == START {
		userRequest[recipientID] = []string{"START"}
		// if !contains(userRequest[recipientID], "START") {
		// 	userRequest[recipientID] = append(userRequest[recipientID], "START")
		// }
		var length = len(allCoins)
		var i int = 0
		for i < length && i < length/3*3 {
			buttons := buttonPostback(allCoins[i : i+3])
			log.Println("< <COINS> > ", allCoins[i:i+3])
			popUpAllCoinButtons(context.TODO(), "Get all coins", recipientID, buttons)
			i += 3
		}

		log.Println("LEngth --> ", i, " ", length)
		buttons := buttonPostback(allCoins[i:length])
		return popUpAllCoinButtons(context.TODO(), "Get all coins", recipientID, buttons)

	} else if msgText == GET_ALL_COINS {
		var length = len(allCoins)
		var i int = 0
		for i < length && i < length/3*3 {
			buttons := buttonPostback(allCoins[i : i+3])
			log.Println("< <COINS> > ", allCoins[i:i+3])
			popUpAllCoinButtons(context.TODO(), "Get all coins", recipientID, buttons)
			i += 3
		}

		buttons := buttonPostback(allCoins[i:length])
		popUpAllCoinButtons(context.TODO(), "Get all coins", recipientID, buttons)
		return nil
	} else if strings.Contains(msgText, GET_FOLLOWED_COINS) {
		var follow, ok = api.GetFollowCoins(recipientID, MESSENGER)
		if ok {
			var length = len(follow)
			var i int = 0
			for i < length && i < length/3*3 {
				buttons := buttonPostback(follow[i : i+3])
				log.Println("< FOLLOW > ", follow[i:i+3])
				popUpAllCoinButtons(context.TODO(), "You have followed these coins", recipientID, buttons)
				i += 3
			}

			buttons := buttonPostback(follow[i:length])
			popUpAllCoinButtons(context.TODO(), "You have followed these coins", recipientID, buttons)

			followedCoins[recipientID] = follow
			return nil
		}
	} else if strings.Contains(msgText, SET_TIME) {
		//80hours10mins
		//^[0-9]:[]
		re, _ := regexp.Compile(SET_TIME + " (.*)")
		submatch := re.FindSubmatch([]byte(msgText))
		userRequest[recipientID] = append(userRequest[recipientID], string(submatch[1]))
		responseText = SET_TIME + " successfully. Wanna save details, please type \"end\" "
	} else if strings.Contains(msgText, SET_BOUND) {
		//upper 3000 lower 200
		re, _ := regexp.Compile(SET_BOUND + " (.*) lower (.*)")
		submatch := re.FindSubmatch([]byte(msgText))
		for _, v := range submatch {
			log.Println(string(v))
		}

		upper := string(submatch[1])
		lower := string(submatch[2])

		userRequest[recipientID] = append(userRequest[recipientID], upper, lower)
		// set upper and lower at the same time

		responseText = "set bounds successfully, please " + SET_TIME + ". Eg: " + SET_TIME + " 10:30AM"
	} else if strings.Contains(msgText, UPDATE_BOUND) {
		re, _ := regexp.Compile(UPDATE_BOUND + " (.*) lower (.*)")
		submatch := re.FindSubmatch([]byte(msgText))
		for _, v := range submatch {
			log.Println(string(v))
		}
		upper := string(submatch[1])
		lower := string(submatch[2])

		userRequest[recipientID] = append(userRequest[recipientID], upper, lower)

		responseText = "Finish update. Type " + END
	} else if msgText == END {
		// giữ phần tử từ 0 -> n - 1  [:n]
		responseText = "end setting for " + userRequest[recipientID][1] + ". Check your followed list by using command <" + GET_FOLLOWED_COINS + ">"
		log.Println("BEFORE -> ", userRequest[recipientID])
		coin_code := userRequest[recipientID][1]
		upper, _ := strconv.ParseFloat(userRequest[recipientID][2], 64)
		lower, _ := strconv.ParseFloat(userRequest[recipientID][3], 64)

		user := api.User{Id: recipientID, Name: recipientID, Platform: MESSENGER}

		// save data
		if api.SetFollowCoin(user, coin_code) {
			log.Println("Follow coin successfully ")
		}
		api.SetLimit(user, api.Limit{Code: coin_code, IsUpper: true, Rate: upper})
		api.SetLimit(user, api.Limit{Code: coin_code, IsUpper: false, Rate: lower})

		if len(userRequest[recipientID]) == 5 {
			time := userRequest[recipientID][4]
			api.SetTime(user, time)
		}

		//delete data
		userRequest[recipientID] = userRequest[recipientID][:1]
		followedCoins[recipientID], _ = api.GetFollowCoins(recipientID, MESSENGER)
		log.Println("AFTER -> ", followedCoins[recipientID])
		// responseText = "helpppp"

	} else {
		responseText = "What can I do for you?"
	}
	return Respond(context.TODO(), recipientID, responseText)
}

func HandlePostBack(recipientID string, payload string) error {
	var responseText string
	if entry, ok := userRequest[recipientID]; ok {
		entry = append(entry, payload)
		userRequest[recipientID] = entry
	} else {
		temp := []string{payload}
		userRequest[recipientID] = temp
	}

	if contains(followedCoins[recipientID], payload) {
		user := api.User{Id: recipientID, Name: recipientID, Platform: MESSENGER}
		info, _ := api.GetLimit(user, payload, true)
		log.Println(info)

		command := []string{UPDATE, DELETE}
		buttons := buttonPostback(command)
		return popUpAllCoinButtons(context.TODO(), "Update or Delete "+userRequest[recipientID][1], recipientID, buttons)
	} else if payload == UPDATE {
		responseText = "Please update bounds. Eg:update upper 23 lower 10"

	} else if payload == DELETE {
		responseText = "Coin " + userRequest[recipientID][1] + " has been deleted. End setting!"

		user := api.User{Id: recipientID, Name: recipientID, Platform: MESSENGER}
		api.DeleteFollowCoin(user, userRequest[recipientID][1])
		userRequest[recipientID] = userRequest[recipientID][:1]
	} else {
		responseText = "Please set bounds. Eg:set upper 3000 lower 200"
	}

	log.Println("user Req", userRequest[recipientID])

	return Respond(context.TODO(), recipientID, responseText)
}

func buttonPostback(buttonList []string) AttachmentButtons {
	var buttons AttachmentButtons
	for i := 0; i < len(buttonList); i++ {
		b := AttachmentButton{
			Type:    "postback",
			Title:   buttonList[i],
			Payload: buttonList[i],
		}
		buttons = append(buttons, b)
	}

	return buttons
}
