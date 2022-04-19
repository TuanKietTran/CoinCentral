package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Coin struct {
	coinid  int
	name    string
	command []string
	price   float32
}

//Getting start:
//1.List all availablle bitcoins
//2.List all followed bitcoins
//3.List help command
const list_bitcoins_command = "List bitcoins"
const list_followed_bitcoins_command = "List followed bitcoins"
const help_command = "help_command"
const follow_coin = "follow"
const select_coin = "select"

var gettingStartKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(list_bitcoins_command, list_bitcoins_command),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(list_followed_bitcoins_command, list_followed_bitcoins_command),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(help_command, help_command),
	),
)

//List all bitcoins: User types
//follow <bitcoin>
// var coinActionKeyboard = tgbotapi.NewInlineKeyboardMarkup(
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
// 		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
// 		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
// 	),
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
// 		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
// 		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
// 	),
// )

// Cannot click, for guide only
// var gettingStartKeyboard = tgbotapi.NewInlineKeyboardMarkup(
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
// 		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
// 		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
// 	),
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
// 		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
// 		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
// 	),
// )

func main() {
	bot, err := tgbotapi.NewBotAPI("5155203320:AAHV4ZhaXl9ByKJtGNKDLIwnsxZ1zWQjAzc")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	//5043456646
	// msg := tgbotapi.NewMessage(5043456646, "Hello world")
	// bot.Send(msg)

	for update := range updates {
		var response string
		var keyboard tgbotapi.InlineKeyboardMarkup
		var chatid int64
		if update.Message != nil {
			log.Printf("------------------------------\n")
			log.Printf("Message [%s] %s", update.Message.From.UserName, update.Message.Text)
			log.Printf("\n-------------------------------\n")
			response, keyboard = get_response(update.Message.Text) // String.
			chatid = update.Message.Chat.ID

		} else {
			log.Printf("------------------------------\n")
			log.Printf("Callback [%s]", update.CallbackQuery.Data)
			log.Printf("\n-------------------------------\n")
			response, keyboard = get_response(update.CallbackQuery.Data)
			chatid = update.CallbackQuery.Message.Chat.ID
		}
		msg := tgbotapi.NewMessage(chatid, response)
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	}

}

func get_response(msg string) (string, tgbotapi.InlineKeyboardMarkup) {
	response := ""
	var keyboard tgbotapi.InlineKeyboardMarkup
	//Getting start
	switch msg {
	case "\\start":
		response = "Bitcoin service, how can I help you?"
		keyboard = gettingStartKeyBoard
		break
	case "\\help_command":
		response = "List all command"
		keyboard = gettingStartKeyBoard
		break
	case list_bitcoins_command:
		response = "Which bitcoin you want to follow\n"
		bitcoins := get_fake_list_bitcoin()
		keyboard = get_bitcoin_keyboards(bitcoins, follow_coin)
		break
	case list_followed_bitcoins_command:
		response = "Select coins for more detail!"
		bitcoins := get_fake_list_followed_bitcoin()
		keyboard = get_bitcoin_keyboards(bitcoins, select_coin)
		break
	default:
		if strings.Contains(msg, follow_coin) {
			//Extract coin id
			coinid, err := strconv.ParseInt(msg[len(follow_coin):], 10, 32)
			if err != nil {
				log.Panic(err)
			}
			log.Printf("Coinid: %d", coinid)
			bitcoins := get_fake_list_bitcoin()
			var coin Coin
			for i := 0; i < len(bitcoins); i++ {
				if int(coinid) == bitcoins[i].coinid {
					coin = bitcoins[i]
				}
			}
			response = get_bitcoin_response(coin)
			keyboard = gettingStartKeyBoard
		} else if strings.Contains(msg, select_coin) {
			//TODO: List coind detail

		}

		break
	}
	return response, keyboard
}

func get_bitcoin_response(coin Coin) string {
	response := "name: " + coin.name + "\nprice: " + fmt.Sprintf("%f", coin.price) + "\nCommands:\n"
	for i := 0; i < len(coin.command); i++ {
		response = response + "[" + strconv.Itoa(i) + "]: " + coin.command[i] + "\n"
	}
	return response
}

func start() string {
	return "Hello, this is coinbot"
}

func help() string {
	return "\\Hello: Say hi to Coinbot\n\\select: Return list of avaiable bitcoins"
}

func get_bitcoin_keyboards(bitcoins []Coin, command string) tgbotapi.InlineKeyboardMarkup {
	var keybuttons = []tgbotapi.InlineKeyboardButton{}
	for i := 0; i < len(bitcoins); i++ {
		keybuttons = append(keybuttons,
			tgbotapi.NewInlineKeyboardButtonData(bitcoins[i].name, command+strconv.Itoa(bitcoins[i].coinid)),
		)
	}
	command_buttons := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Return", "\\start"),
	)
	return tgbotapi.NewInlineKeyboardMarkup(
		keybuttons,
		command_buttons)
}

func get_fake_list_bitcoin() []Coin {
	var bitcoins = []Coin{
		{0, "Acoin", []string{"Upper bound 1", "Upper bound 2", "Lower bound 3"}, 200},
		{1, "Bcoin", []string{"Upper bound 1", "Upper bound 2", "Lower bound 3"}, 300},
		{2, "Ccoin", []string{"Upper bound 1", "Upper bound 2", "Lower bound 3"}, 400},
		{3, "Dcoin", []string{"Upper bound 1", "Upper bound 2", "Lower bound 3"}, 500},
		{4, "Ecoin", []string{"Upper bound 1", "Upper bound 2", "Lower bound 3"}, 600},
		{5, "Fcoin", []string{"Upper bound 1", "Upper bound 2", "Lower bound 3"}, 900},
		{6, "Gcoin", []string{"Upper bound 1", "Upper bound 2", "Lower bound 3"}, 1000},
	}
	return bitcoins
}

func get_fake_list_followed_bitcoin() []Coin {
	var bitcoins = []Coin{
		{0, "Acoin", []string{"Upper bound 1", "Upper bound 2", "Lower bound 3"}, 200},
		{1, "Bcoin", []string{"Upper bound 1", "Upper bound 2", "Lower bound 3"}, 300},
		{2, "Ccoin", []string{"Upper bound 1", "Upper bound 2", "Lower bound 3"}, 400},
	}
	return bitcoins
}
