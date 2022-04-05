package models

import (
	"time"
)

type CoinMeta struct {
	Code              string
	Name              string
	Symbol            string
	Rank              int
	Age               int
	Color             string
	Png32             string
	Png64             string
	Webp32            string
	Webp64            string
	Exchanges         int
	Markets           int
	Pairs             int
	AllTimeHighUSD    float64
	CirculatingSupply int
	TotalSupply       int
	MaxSupply         int
}

type CoinRate struct {
	Code       string
	Rate       float64
	Volume     int
	Cap        int
	InsertTime time.Time
}
