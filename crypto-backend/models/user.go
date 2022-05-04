package models

type User struct {
	UserId      string  		`bson:"_id" json:"userId"`
	Name        string  		`bson:"name" json:"name"`
	TelegramId  int     		`bson:"telegramId" json:"telegramId"`
	MessengerId string  		`bson:"messengerId" json:"messengerId"`
	LimitList   []Limit 		`bson:"limitList" json:"limitList"`
	WatchList   Notification 	`bson:"watchList" json:"watchList"`
}

type Limit struct {
	Code    string  `bson:"code" json:"code"`
	IsUpper bool    `bson:"isUpper" json:"isUpper"`
	Rate    float64 `bson:"rate" json:"rate"`
}

type Notification struct {
	Code []string `bson:"code" json:"code"`
	Time []string `bson:"time" json:"time"`
}
