package models

import "time"

type UserId struct {
	Id       string `bson:"id" json:"id"`
	Platform string `bson:"platform" json:"platform"`
}

type User struct {
	Id        string      `bson:"id" json:"id"`
	Platform  string      `bson:"platform" json:"platform"`
	Name      string      `bson:"name" json:"name"`
	LimitList []Limit     `bson:"limitList" json:"limitList"`
	CodeList  []string    `bson:"codeList" json:"codeList"`
	TimeList  []time.Time `bson:"timeList" json:"timeList"`
}

type Limit struct {
	Code    string  `bson:"code" json:"code"`
	IsUpper bool    `bson:"isUpper" json:"isUpper"`
	Rate    float64 `bson:"rate" json:"rate"`
}
