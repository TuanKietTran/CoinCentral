package api

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
