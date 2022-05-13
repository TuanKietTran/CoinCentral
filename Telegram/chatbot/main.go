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
	URL   = "https://9f89-14-169-108-130.ap.ngrok.io"
	PORT  = ":8089"
)

//GLOBAL VARIABLE
var bot *tgbotapi.BotAPI
var err error
var user_state_list (map[string](object.State))
var coins []object.Coin

// COMMAND LIST
const (
	START                   = "\\start"
	LIST_ALL_COINS          = "GET ALL COINS"
	LIST_ALL_FOLLOWED_COINS = "GET ALL FOLLOW COINS"
	HELP                    = "HELP"
	UNFOLLOW_COIN_ACTION    = "UNFOLLOW"
	FOLLOW_COIN_ACTION      = "FOLLOW"
	SELECT_COIN_ACTION      = "SELECT"
	TIME_ADD                = "TIME_ADD"
	ADD_TIME                = "ADD TIME"
	RETURN                  = "RETURN"
	EDIT                    = "EDIT"
	UPPER                   = "UPPER"
	LOWER                   = "LOWER"
	TIME                    = "TIME"
	EDIT_TIME               = "EDIT_TIME"
	UPDATE_TIME             = "UPDATE TIME"
	DEFAUTL_TIME            = "5:30PM"
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
		tgbotapi.NewInlineKeyboardButtonData(TIME, TIME),
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

//Callback funciton from webhook:
var limitMsgPostBacker = func(limitMsg object.WebhookLimitMsg) {
	if bot != nil {
		user := limitMsg.UserId
		limit := limitMsg.Limit

		fmt.Println("=============================")
		fmt.Println(user.Id)
		fmt.Println(limit)

		state := "No"
		if limit.IsUpper {
			state = "Yes"
		}

		reply := "Coin: " + limit.Code +
			"\nisUpper: " + state +
			"\nrate:" + fmt.Sprint(limit.Rate) +
			"\n"

		id, err := strconv.ParseInt(user.Id, 10, 64)
		if err != nil {
			panic(err)
		}
		msg := tgbotapi.NewMessage(id, reply)
		bot.Send(msg)
	}
}

var timeMsgPostBacker = func(timeMsg object.WebhookTimeMsg) {
	if bot != nil {
		userid := timeMsg.UserId
		coins := timeMsg.Coins

		var reply string
		for i, coin := range coins {
			reply += fmt.Sprint(i) + ". Coin: " + coin.Code +
				"\n  rate:" + fmt.Sprint(coin.Rate) +
				"\n"
		}

		fmt.Println(reply)

		id, err := strconv.ParseInt(userid, 10, 64)
		if err != nil {
			panic(err)
		}
		msg := tgbotapi.NewMessage(id, reply)
		bot.Send(msg)
	}
}

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
	user_state_list = make(map[string]object.State)
	//Clone data:
	coins = get_all_bitcoins()

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
		add_states(chatid, START)                         //Save current state
		response = "Bitcoin service, how can I help you?" //Get response
		keyboard = gettingStartKeyBoard                   //Get keyboard
		isUsedKeyBoard = true
		success := api.CreateUser(user)
		if !success {
			fmt.Println("User may already exist!.")
		}
		//Download data:
		follow_coins := get_follow_bitcoins(user)
		times, _ := api.GetTime(user.Id, user.Platform)
		user_state_list[chatid] = object.State{[]string{START}, follow_coins, times}
		break
	case HELP: //STEP 1.3
		//Remove all previous state (except start)
		response = "List all command"
		keyboard = gettingStartKeyBoard
		isUsedKeyBoard = true
		break
	case LIST_ALL_COINS: //STEP 1.1:
		add_states(chatid, LIST_ALL_COINS)
		response = "Which bitcoin you want to follow\n"
		//Only return non-folow coins
		unfollow_coins := get_list_unfollow_coins(coins, user_state_list[chatid].Coins)
		if len(unfollow_coins) == 0 {
			response = "You have follow all bitcoins we haved!"
			isUsedKeyBoard = false
		} else {
			keyboard = get_bitcoin_keyboards(unfollow_coins, SELECT_COIN_ACTION)
			isUsedKeyBoard = true
		}
		break
	case LIST_ALL_FOLLOWED_COINS: //STEP 1.2
		add_states(chatid, LIST_ALL_FOLLOWED_COINS)
		response = "Select coins for more detail!"
		if len(user_state_list[chatid].Coins) == 0 {
			isUsedKeyBoard = false
		} else {
			keyboard = get_bitcoin_keyboards(user_state_list[chatid].Coins, SELECT_COIN_ACTION) //Get keyboard show all followed coins
			isUsedKeyBoard = true
		}
		break
	case TIME:
		add_states(chatid, TIME)
		response = "Press button that you want to edit time!"
		keyboard = getTimeKeyBoard(user_state_list[chatid].Times, EDIT_TIME)
		isUsedKeyBoard = true
		break
	case FOLLOW_COIN_ACTION: //STEP3 (FROM STEP LIST ALL => SELECT => FOLLOW)
		add_states(chatid, FOLLOW_COIN_ACTION)
		response = "Input your upper bound, " + UPPER + " <Rate>" //Get response
		isUsedKeyBoard = false
		//No keyboard is used, as user's going to type a command
		break
	case EDIT: //STEP 3: (FROM STEP LIST ALL FOLLOWED COINS => SELECT => EDIT)
		add_states(chatid, EDIT)
		response = "Which action that you want to do." //Get response
		keyboard = upperLowerUnfollowedKeyBoard        //Get keyboard: User choose upper,lower, or unfollowed coins
		isUsedKeyBoard = true
		break
	case UNFOLLOW_COIN_ACTION: //STEP 4.1 (FOLLOW STEP 3 RIGHT ABOVE)
		//Get coin_code:
		var coin_code string
		for i, state := range user_state_list[chatid].Step {
			if strings.Contains(state, SELECT_COIN_ACTION) {
				coin_code = state[len(SELECT_COIN_ACTION):]
				remove_follow_coin_at(chatid, i)
				break
			}
		}
		reset_states(chatid)            //Restore state to current first state
		keyboard = gettingStartKeyBoard //Get keyboard
		isUsedKeyBoard = true
		//TODO: Unnfollow coin
		success := api.DeleteFollowCoin(user, coin_code)
		if !success {
			response = "Unfollow UNsuccessful!" //Get response
		} else {
			response = "Unfollow successful!"
		}

		break
	case UPPER: //Step 4.2  (FOLLOW STEP 3 RIGHT ABOVE)
		add_states(chatid, UPPER)
		response = "Input your upper bound," + EDIT + " " + UPPER + " <Rate>" //get response
		isUsedKeyBoard = false                                                //User's going to type command
		break
	case LOWER: //Step 4.3 (FOLLOW STEP 3 RIGHT ABOVE)
		add_states(chatid, LOWER)
		response = "Input your lower bound, " + EDIT + " " + LOWER + " <Rate>"
		isUsedKeyBoard = false
		break
	case RETURN: //Abort action, return to first state
		reset_states(chatid)
		response = "Bitcoin service, how can I help you?"
		isUsedKeyBoard = true
		keyboard = gettingStartKeyBoard
		break
	case TIME_ADD:
		if len(user_state_list[chatid].Times) < 5 {
			add_states(chatid, TIME_ADD)
			isUsedKeyBoard = false
			response = "Please type: ADD TIME 00:00AM for ADD time"
		} else {
			response = "You can only set 5 times!"
			keyboard = getTimeKeyBoard(user_state_list[chatid].Times, EDIT_TIME)
			isUsedKeyBoard = true
		}
		break
	default:
		if strings.Contains(msg, SELECT_COIN_ACTION) { //STEP 2: FOR BOTH LIST ALL AND FOLLWED COINS
			add_states(chatid, msg)                   //Save ACTION + COIN CODE
			var selected_coin object.Coin             //Coin that selected by user later
			coincode := msg[len(SELECT_COIN_ACTION):] //msg at this state is SELECT<COIN CODE>, so we ignore SELECT to got code.
			for _, coin := range coins {
				if coincode == coin.Code {
					selected_coin = coin
					break
				}
			}
			response = "Code: " + selected_coin.Code + //Get response
				"\nName: " + selected_coin.Name +
				"\nRate: " + fmt.Sprintf("%f", selected_coin.Rate)
			if user_state_list[chatid].Step[1] == LIST_ALL_COINS { //If user get list of all nodes, return "followSelected" keyboard
				keyboard = followSelectedCoinKeyboard //It's return 2 button: Follow (if user want follow coin) and return
				isUsedKeyBoard = true
			} else {
				keyboard = selecteFollowedCoinKeyboard //Else, return "selectfollowed" coins.
				isUsedKeyBoard = true                  //It's return Edit (wether user need update coins limit, or unfollowed coins)
			}
			break
		} else if strings.Contains(msg, EDIT) && strings.Contains(msg, UPPER) { //STEP 4: LIST FOLLOWED COINS -> EDIT ->
			var coin_code string
			var upper float64
			for _, state := range user_state_list[chatid].Step {
				upper_state := strings.ToUpper(state)
				if strings.Contains(upper_state, SELECT_COIN_ACTION) {
					coin_code = state[len(SELECT_COIN_ACTION):]
				}
				//Eliminate space:
			}
			pre_state := strings.ReplaceAll(msg, " ", "")
			upper, _ = strconv.ParseFloat(pre_state[(len(EDIT)+len(UPPER)):], 64)
			fmt.Println("Upper: ", upper)
			keyboard = gettingStartKeyBoard
			isUsedKeyBoard = true
			reset_states(chatid)

			//SAVE NEW LIMIT UPPER BOUND
			limit := object.Limit{coin_code, true, upper}
			success := api.UpdateLimit(user, limit)
			if success {
				response = "Update new limit success!"
			} else {
				response = "Update limit, retry latter"
			}
		} else if strings.Contains(msg, EDIT) && strings.Contains(msg, LOWER) {
			var coin_code string
			var lower float64
			for _, state := range user_state_list[chatid].Step {
				upper_state := strings.ToUpper(state)
				if strings.Contains(upper_state, SELECT_COIN_ACTION) {
					coin_code = state[len(SELECT_COIN_ACTION):]
				}
			}
			pre_state := strings.ReplaceAll(msg, " ", "")
			lower, _ = strconv.ParseFloat(pre_state[(len(EDIT)+len(LOWER)):], 64)
			fmt.Println("Lower: ", lower)
			keyboard = gettingStartKeyBoard
			isUsedKeyBoard = true
			reset_states(chatid)

			//TODO:: SAVE NEW LIMIT  LOWER BOUND
			limit := object.Limit{coin_code, false, lower}
			success := api.UpdateLimit(user, limit)
			if success {
				response = "Update new limit success!"
			} else {
				response = "Update limit, retry latter"
			}
		} else if strings.Contains(msg, UPPER) { //STEP 4: FOR LIST ALL COINS -> FOLLOW -> UPPER
			add_states(chatid, msg) //Save Action and upper bound
			response = "Input your lower bound, " + LOWER + " <Rate>"
			isUsedKeyBoard = false //user's going to type a command
		} else if strings.Contains(msg, LOWER) { //STEP 5: FOR LIST ALL COINS -> FOLLOW -> LOWER
			//Extract coin code, upper bound and lower bound
			var coin_code string
			var upper float64
			var lower float64
			for _, state := range user_state_list[chatid].Step {
				upper_state := strings.ToUpper(state)
				if strings.Contains(upper_state, SELECT_COIN_ACTION) {
					coin_code = state[len(SELECT_COIN_ACTION):]
				}
				if strings.Contains(upper_state, UPPER) {
					//Eliminate space:
					pre_state := strings.ReplaceAll(state, " ", "")
					upper, _ = strconv.ParseFloat(pre_state[len(UPPER):], 64)
				}
			}
			//Eliminate space:
			pre_state := strings.ReplaceAll(msg, " ", "")
			lower, _ = strconv.ParseFloat(pre_state[len(LOWER):], 64)
			success := api.SetFollowCoin(user, coin_code)
			if !success {
				response = "System can not register coin for user!"
			} else {
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
				//Add selected coin for current user cache
				for _, coin := range coins {
					if coin.Code == coin_code {
						entry, _ := user_state_list[chatid]
						entry.Coins = append(entry.Coins, coin)
						user_state_list[chatid] = entry
						break
					}
				}
			}

			keyboard = gettingStartKeyBoard
			reset_states(chatid)
			isUsedKeyBoard = true
		} else if strings.Contains(msg, EDIT_TIME) {
			add_states(chatid, msg)
			response = "Command: UPDATE TIME 00:00AM for update time"
			isUsedKeyBoard = false
		} else if strings.Contains(msg, UPDATE_TIME) {
			var index int64
			var time string
			for _, state := range user_state_list[chatid].Step {
				state = strings.ToUpper(state)
				if strings.Contains(state, EDIT_TIME) {
					index, _ = strconv.ParseInt(state[len(EDIT_TIME):], 10, 32)
				}
			}
			pre_state := strings.ReplaceAll(msg, " ", "")
			time = pre_state[(len(UPDATE_TIME) - 1):]

			fmt.Println("idx ", index)
			fmt.Println("Time ", time)
			if index < int64(len(user_state_list[chatid].Times)) {
				old_time := user_state_list[chatid].Times[index]
				api.DeleteTime(user, old_time)
				user_state_list[chatid].Times[index] = time
			}
			success := api.SetTime(user, time)
			if success {
				response = "Update successful!"
			} else {
				response = "Update fail!"
			}
			isUsedKeyBoard = true
			keyboard = gettingStartKeyBoard
			reset_states(chatid)

		} else if strings.Contains(msg, ADD_TIME) {
			pre_state := strings.ReplaceAll(msg, " ", "")
			time := pre_state[(len(ADD_TIME) - 1):]
			fmt.Println("time: ", time)
			success := api.SetTime(user, time)
			if success {
				response = "Update successful!"
			} else {
				response = "Update fail!"
			}

			if entry, ok := user_state_list[chatid]; ok {
				entry.Times = append(entry.Times, time)
				user_state_list[chatid] = entry
			}

			isUsedKeyBoard = true
			keyboard = gettingStartKeyBoard
			reset_states(chatid)
		}
	}
	fmt.Println("=============================")
	fmt.Println(user_state_list)
	fmt.Println("===============================")
	return response, keyboard, isUsedKeyBoard
}

func help() string {
	return "\\Hello: Say hi to Coinbot\n\\select: Return list of avaiable bitcoins"
}

func get_list_unfollow_coins(full_list []object.Coin,
	follow_list []object.Coin) []object.Coin {
	var coins []object.Coin
	var isfollow bool
	for _, coin := range full_list {
		isfollow = false
		for _, follow_coin := range follow_list {
			if coin.Code == follow_coin.Code {
				isfollow = true
				break
			}
		}
		if !isfollow {
			coins = append(coins, coin)
		}
	}
	return coins
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

func get_all_bitcoins() []object.Coin {
	//Get list bitcoin object from list of bitcoins name
	coinCodes := api.GetAllCoins()
	var coinList []object.Coin
	for _, coin := range coinCodes {
		rcoin, _ := api.GetCoin(coin)
		coinList = append(coinList, rcoin)
	}
	return coinList
}

func get_follow_bitcoins(user object.User) []object.Coin {
	coinCodes, _ := api.GetFollowCoins(user.Id, user.Platform)
	var coinList []object.Coin
	for _, coin := range coinCodes {
		rcoin, _ := api.GetCoin(coin)
		coinList = append(coinList, rcoin)
	}
	fmt.Println("Coin codes: ", coinCodes)
	return coinList
}

func add_states(key string, step string) {
	if entry, ok := user_state_list[key]; ok {
		entry.Step = append(entry.Step, step)
		user_state_list[key] = entry
	} else {
		user_state_list[key] = object.State{}
	}
}

func reset_states(key string) {
	if entry, ok := user_state_list[key]; ok {
		entry.Step = entry.Step[:1]
		user_state_list[key] = entry
	}
}

func remove_follow_coin_at(key string, i int) {
	if entry, ok := user_state_list[key]; ok {
		entry.Step = append(entry.Step[:i], entry.Step[i+1:]...)
		user_state_list[key] = entry
	}
}

//WEB HOOK CLIENT
func handling_http(bot *tgbotapi.BotAPI) {
	//Sending url for webhook: Only for ngrok
	platform := "telegram"
	request := object.WebhookRequest{URL, platform}
	api.CreateWebhookRequest(request)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Println("Hello world\n") })
	http.HandleFunc("/ping", PingHandler)
	http.HandleFunc("/limits", LimitsHandler)
	http.HandleFunc("/times", TimesHandler)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getTimeKeyBoard(times []string, command string) tgbotapi.InlineKeyboardMarkup {
	//Return keyboard
	var keyrows [][]tgbotapi.InlineKeyboardButton
	for _, time := range times {
		keybuttons := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				time, command+time))
		keyrows = append(keyrows, keybuttons)
	}
	add_buttons := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("ADD TIME", TIME_ADD),
	)
	return_buttons := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(RETURN, RETURN),
	)
	keyrows = append(keyrows, add_buttons, return_buttons)
	return tgbotapi.NewInlineKeyboardMarkup(
		keyrows...)
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
