package api

// type UserMapIDKey = map[string]UserData
// type CoinMapNameKey = map[string]CoinData

// type (
// 	UserData struct {
// 		State string //full, pend
// 		Data  RequestData
// 	}

// 	RequestData struct {
// 		FollowedCoinList CoinMapNameKey
// 		TimeInterval     int //mins
// 	}

// 	CoinData struct {
// 		Name       string
// 		Upperbound int
// 		Lowerbound int
// 	}
// )

//STEP 2: LIST ALL COINS , FOLLOW COINS, CURRENT <TIME>, HELP COMMAND
//STEP 3: FOLLOW, RETURN TO START
//
// sender_id: [GET ALL COINS","<coin name>", "<upper>", "<lower>", "<time>", "END"]
type UserRequest = map[string]([]string)

type User struct {
	Id       string `json:"id"`
	Platform string `json:"platform"`
	Name     string `json:"name"`
}

type Coin struct {
	Code string  `json:"code"`
	Name string  `json:"name"`
	Rate float64 `json:"rate"`
}

type Limit struct {
	Code    string  `json:"code"`
	IsUpper bool    `json:"isUpper"`
	Rate    float64 `json:"rate"`
}

type WebhookRequest struct {
	CallbackUrl string `json:"callbackUrl"`
	Platform    string `json:"platform"`
}

type WebhookLimitMsg struct {
	UserId User  `json:"userId"`
	Limit  Limit `json:"limit"`
}
