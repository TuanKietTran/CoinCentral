module github.com/TuanKietTran/CoinCentral/telebot

go 1.18

replace github.com/TuanKietTran/CoinCentral/api => ../api

replace github.com/TuanKietTran/CoinCentral/object => ../object

require (
	github.com/TuanKietTran/CoinCentral/api v0.0.0-00010101000000-000000000000
	github.com/TuanKietTran/CoinCentral/object v0.0.0-00010101000000-000000000000
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
)
