package models

type User struct {
	UserId        string      `bson:"_id" json:"userId"`
	Name          string      `bson:"name" json:"name"`
	ThresholdList []Threshold `bson:"thresholdList" json:"-"`
}

type Threshold struct {
	Code  string  `bson:"code" json:"code"`
	Upper float64 `bson:"upper,omitempty" json:"upper"`
	Lower float64 `bson:"lower,omitempty" json:"lower"`
}

type Period struct {
	Code string `bson:"code" json:"code"`
	Time int64  `bson:"time"`
}
