package webhook

type URL struct {
	TelegramBotConnected bool
	TelegramCallbackUrl  string

	MessengerBotConnected bool
	MessengerCallbackUrl  string
}
