package models

type Coin struct {
	Code string  `bson:"code" json:"code"`
	Name string  `bson:"name" json:"name"`
	Rate float64 `bson:"rate" json:"rate"`
}
