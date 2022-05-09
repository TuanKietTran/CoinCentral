package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/TuanKietTran/CoinCentral/api"
	"github.com/TuanKietTran/CoinCentral/object"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	STEP0 = 0
	STEP1 = 1
	STEP2 = 2
	STEP3 = 3
	STEP4 = 4
	STEP5 = 5
)

// COMMAND LIST

const (
	START                   = "\\start"
	LIST_ALL_COINS          = "GET ALL COINS"
	LIST_ALL_FOLLOWED_COINS = "GET ALL FOLLOW COINS"
	HELP                    = "HELP"
	UNFOLLOW_COIN_ACTION    = "UNFOLLOW"
	FOLLOW_COIN_ACTION      = "FOLLOW"
	SELECT_COIN_ACTION      = "SELECT"
	RETURN                  = "RETURN"
	EDIT                    = "EDIT"
	UPPER                   = "UPPER"
	LOWER                   = "LOWER"
)

//KEY BOARD LIST
var gettingStartKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(LIST_ALL_COINS, LIST_ALL_COINS),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(LIST_ALL_FOLLOWED_COINS, LIST_ALL_FOLLOWED_COINS),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(HELP, HELP),
	),
)

var followSelectedCoinKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(FOLLOW_COIN_ACTION, FOLLOW_COIN_ACTION),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(RETURN, RETURN),
	),
)

var selectCoinKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(SELECT_COIN_ACTION, SELECT_COIN_ACTION),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(RETURN, RETURN),
	),
)

var upperLowerUnfollowedKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(UPPER, UPPER),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(LOWER, LOWER),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(UNFOLLOW_COIN_ACTION, UNFOLLOW_COIN_ACTION),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(RETURN, RETURN),
	),
)

var selecteFollowedCoinKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(EDIT, EDIT),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(RETURN, RETURN),
	),
)

var returnKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(RETURN, RETURN),
	),
)

//GLOBAL VARIABLE
var bot *tgbotapi.BotAPI
var err error
var user_state_list (map[string]([]string))
var coins []object.Coin

func main() {
	//Initialize App:
	bot, err = tgbotapi.NewBotAPI("5155203320:AAHV4ZhaXl9ByKJtGNKDLIwnsxZ1zWQjAzc")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	user_state_list = make(map[string][]string)

	//Clone data:
	coins = get_bitcoins()
	fmt.Println(coins)

	//Turn on webhook
	go handling_http(bot)

	for update := range updates {
		var response string
		var keyboard tgbotapi.InlineKeyboardMarkup
		var chatid int64
		var isUsedKeyBoard bool
		if update.Message != nil {
			//When user command
			// log.Printf("------------------------------\n")
			// log.Println("Message: ", update)
			// log.Printf("\n-------------------------------\n")
			chatid = update.Message.Chat.ID
			user := object.User{strconv.FormatInt(update.Message.Chat.ID, 10), "telegram", update.Message.Chat.FirstName}
			response, keyboard, isUsedKeyBoard = get_response(user, update.Message.Text) // String.

		} else {
			//When user press keyboard
			// log.Printf("------------------------------\n")
			// log.Printf("Callback [%s]", update.CallbackQuery.Data)
			// log.Printf("\n-------------------------------\n")
			chatid = update.CallbackQuery.Message.Chat.ID
			user := object.User{strconv.FormatInt(update.CallbackQuery.Message.Chat.ID, 10),
				"telegram",
				update.CallbackQuery.Message.Chat.FirstName}
			response, keyboard, isUsedKeyBoard = get_response(user, update.CallbackQuery.Data)
		}
		msg := tgbotapi.NewMessage(chatid, response)
		if isUsedKeyBoard {
			msg.ReplyMarkup = keyboard
		}
		bot.Send(msg)
	}
}

func get_response(user object.User, msg string) (string, tgbotapi.InlineKeyboardMarkup, bool) {
	/*
		args chatid: Telegram's UserID
		args msg: Message receive from User
		return response, keyboard, isUsedKeyBoard
	*/
	chatid := user.Id
	var response string
	var keyboard tgbotapi.InlineKeyboardMarkup // Keyboard
	var isUsedKeyBoard bool                    //Whether to return keyboard, or let user make a command

	switch msg {
	case START: //FIRST 0
		user_state_list[chatid] = append(user_state_list[chatid], START) //Save current state
		response = "Bitcoin service, how can I help you?"                //Get response
		keyboard = gettingStartKeyBoard                                  //Get keyboard
		isUsedKeyBoard = true
		success := api.CreateUser(user)
		if !success {
			fmt.Println("User may already exist!.")
		}
		break
	case HELP: //STEP 1.3
		//Remove all previous state (except start)
		response = "List all command"
		keyboard = gettingStartKeyBoard
		isUsedKeyBoard = true
		break
	case LIST_ALL_COINS: //STEP 1.1:
		user_state_list[chatid] = append(user_state_list[chatid], LIST_ALL_COINS) //Save current state
		response = "Which bitcoin you want to follow\n"                           //Get response
		keyboard = get_bitcoin_keyboards(coins, SELECT_COIN_ACTION)               //Get keyboard showall available coins
		isUsedKeyBoard = true
		break
	case LIST_ALL_FOLLOWED_COINS: //STEP 1.2
		user_state_list[chatid] = append(user_state_list[chatid], LIST_ALL_FOLLOWED_COINS) //Save current state
		response = "Select coins for more detail!"                                         //Get response
		//TODO: Add list followed bitcoins
		keyboard = get_bitcoin_keyboards(coins, SELECT_COIN_ACTION) //Get keyboard show all followed coins
		isUsedKeyBoard = true
		break
	case FOLLOW_COIN_ACTION: //STEP3 (FROM STEP LIST ALL => SELECT => FOLLOW)
		user_state_list[chatid] = append(user_state_list[chatid], FOLLOW_COIN_ACTION) //Save current state
		response = "Input your upper bound, " + UPPER + " <Rate>"                     //Get response
		isUsedKeyBoard = false                                                        //No keyboard is used, as user's going to type a command
		break
	case EDIT: //STEP 3: (FROM STEP LIST ALL FOLLOWED COINS => SELECT => EDIT)
		user_state_list[chatid] = append(user_state_list[chatid], EDIT) //Save current state
		response = "Which action that you want to do."                  //Get response
		keyboard = upperLowerUnfollowedKeyBoard                         //Get keyboard: User choose upper,lower, or unfollowed coins
		isUsedKeyBoard = true
		break
	case UNFOLLOW_COIN_ACTION: //STEP 4.1 (FOLLOW STEP 3 RIGHT ABOVE)
		//Get coin_code:
		for _, state := range user_state_list[chatid] {
			if strings.Contains(state, SELECT_COIN_ACTION) {
				//coin_code := state[len(SELECT_COIN_ACTION):]
				//TODO: unfollowed code
			}
		}
		user_state_list[chatid] = user_state_list[chatid][:1] //Restore state to current first state
		response = "Unfollow success"                         //Get response
		keyboard = gettingStartKeyBoard                       //Get keyboard
		isUsedKeyBoard = true
		break
	case UPPER: //Step 4.2  (FOLLOW STEP 3 RIGHT ABOVE)
		user_state_list[chatid] = append(user_state_list[chatid], UPPER)      //Save state
		response = "Input your upper bound," + EDIT + " " + UPPER + " <Rate>" //get response
		isUsedKeyBoard = false                                                //User's going to type command
		break
	case LOWER: //Step 4.3 (FOLLOW STEP 3 RIGHT ABOVE)
		user_state_list[chatid] = append(user_state_list[chatid], LOWER)
		response = "Input your lower bound, " + EDIT + " " + LOWER + " <Rate>"
		isUsedKeyBoard = false
		break
	case RETURN: //Abort action, return to first state
		user_state_list[chatid] = user_state_list[chatid][:1]
		response = "Bitcoin service, how can I help you?"
		isUsedKeyBoard = true
		keyboard = gettingStartKeyBoard
		break
	default:
		if strings.Contains(msg, SELECT_COIN_ACTION) { //STEP 2: FOR BOTH LIST ALL AND FOLLWED COINS
			user_state_list[chatid] = append(user_state_list[chatid], msg) //Save ACTION + COIN CODE
			var selected_coin object.Coin                                  //Coin that selected by user later
			coincode := msg[len(SELECT_COIN_ACTION):]                      //msg at this state is SELECT<COIN CODE>, so we ignore SELECT to got code.
			for _, coin := range coins {
				if coincode == coin.Code {
					selected_coin = coin
					break
				}
			}
			response = "Code: " + selected_coin.Code + //Get response
				"\nName: " + selected_coin.Name +
				"\nRate: " + fmt.Sprintf("%f", selected_coin.Rate)
			if user_state_list[chatid][1] == LIST_ALL_COINS { //If user get list of all nodes, return "followSelected" keyboard
				keyboard = followSelectedCoinKeyboard //It's return 2 button: Follow (if user want follow coin) and return
				isUsedKeyBoard = true
			} else {
				keyboard = selecteFollowedCoinKeyboard //Else, return "selectfollowed" coins.
				isUsedKeyBoard = true                  //It's return Edit (wether user need update coins limit, or unfollowed coins)
			}
			break
		} else if strings.Contains(msg, UPPER) { //STEP 4: FOR LIST ALL COINS -> FOLLOW -> UPPER
			user_state_list[chatid] = append(user_state_list[chatid], msg) //Save Action and upper bound
			response = "Input your lower bound, " + LOWER + " <Rate>"
			isUsedKeyBoard = false //user's going to type a command
		} else if strings.Contains(msg, LOWER) { //STEP 5: FOR LIST ALL COINS -> FOLLOW -> LOWER
			//Extract coin code, upper bound and lower bound
			var coin_code string
			var upper float64
			var lower float64
			for _, state := range user_state_list[chatid] {
				upper_state := strings.ToUpper(state)
				if strings.Contains(upper_state, SELECT_COIN_ACTION) {
					coin_code = state[len(SELECT_COIN_ACTION):]
				}
				if strings.Contains(upper_state, UPPER) {
					//Eliminate space:
					pre_state := strings.ReplaceAll(state, " ", "")
					upper, _ = strconv.ParseFloat(pre_state[len(UPPER):], 64)
				}
				if strings.Contains(upper_state, LOWER) {
					//Eliminate space:
					pre_state := strings.ReplaceAll(state, " ", "")
					lower, _ = strconv.ParseFloat(pre_state[len(LOWER):], 64)
				}
			}
			upper_limit := object.Limit{coin_code, true, upper}
			lower_limit := object.Limit{coin_code, false, lower}
			upper_success := api.SetLimit(user, upper_limit)
			lower_success := api.SetLimit(user, lower_limit)
			if !upper_success || !lower_success {
				fmt.Println("Error occur upper: ", upper_success, "lower: ", lower_success)
				response = "Following fails!"
			} else {
				response = "Following success!"
			}
			keyboard = get_bitcoin_keyboards(coins, SELECT_COIN_ACTION)
			user_state_list[chatid] = user_state_list[chatid][:1] //AFTER COMPLETE, RESTORE STATE AND RETURN TO FIRST STEP.
			isUsedKeyBoard = true
			keyboard = gettingStartKeyBoard
			//TODO: SAVE LIMIT
		} else if strings.Contains(msg, EDIT) && strings.Contains(msg, UPPER) { //STEP 4: LIST FOLLOWED COINS -> EDIT ->

			response = "Alter limit success!"
			keyboard = gettingStartKeyBoard
			isUsedKeyBoard = true
			user_state_list[chatid] = user_state_list[chatid][:1] //AFTER EDIT LIMIT, RETURN TO FIRST STEP
			//SAVE NEW LIMIT

		} else if strings.Contains(msg, EDIT) && strings.Contains(msg, LOWER) {
			response = "Alter limit success!"
			keyboard = gettingStartKeyBoard
			isUsedKeyBoard = true
			user_state_list[chatid] = user_state_list[chatid][:1]
		}
	}
	return response, keyboard, isUsedKeyBoard
}

func help() string {
	return "\\Hello: Say hi to Coinbot\n\\select: Return list of avaiable bitcoins"
}

func get_bitcoin_keyboards(bitcoins []object.Coin, command string) tgbotapi.InlineKeyboardMarkup {
	//Return keyboard
	var keyrows [][]tgbotapi.InlineKeyboardButton
	for _, coin := range bitcoins {
		keybuttons := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				coin.Name, command+coin.Code))
		keyrows = append(keyrows, keybuttons)
	}
	command_buttons := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Return", "\\start"),
	)
	keyrows = append(keyrows, command_buttons)
	return tgbotapi.NewInlineKeyboardMarkup(
		keyrows...)
}

func get_bitcoins() []object.Coin {
	//Get list bitcoin object from list of bitcoins name
	coinCodes := api.GetAllCoins()
	var coinList []object.Coin
	for _, coinCode := range coinCodes {
		coin, success := api.GetCoin(coinCode)
		if !success {
			fmt.Println("Could not take api code: " + coinCode)
		}
		coinList = append(coinList, coin)
	}
	return coinList
}

//WEB HOOK CLIENT
func handling_http(bot *tgbotapi.BotAPI) {
	http.HandleFunc("/", Parse)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func Parse(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if len(body) == 0 {
		return
	}

	//Visualize msg
	coinsJson := string(body)
	fmt.Printf(coinsJson)

	var limitMsgs []object.WebhookLimitMsg
	json.Unmarshal(body, &limitMsgs)

	for _, msg := range limitMsgs {
		if msg.UserId.Platform == "telegram" {
			isupper := "No"
			if msg.Limit.IsUpper {
				isupper = "Yes"
			}
			userid, err := strconv.ParseInt(msg.UserId.Id, 10, 64)
			if err != nil {
				fmt.Println("UserID cannot convert to Int")
				continue
			}

			content := "Notice:" +
				"\nCode: " + msg.Limit.Code +
				"\nUpper: " + isupper +
				"\nRate: " + fmt.Sprintf("%f", msg.Limit.Rate)
			msg := tgbotapi.NewMessage(userid, content)
			bot.Send(msg)
		}
	}
}
