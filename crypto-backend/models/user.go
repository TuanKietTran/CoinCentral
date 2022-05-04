package models

type UserId struct {
	Id       string `bson:"id" json:"id"`
	Platform string `bson:"platform" json:"platform"`
}

type User struct {
	Id        string       `bson:"id" json:"id"`
	Platform  string       `bson:"platform" json:"platform"`
	Name      string       `bson:"name" json:"name"`
	LimitList []Limit      `bson:"limitList" json:"limitList"`
	WatchList Notification `bson:"watchList" json:"watchList"`
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
