package models

type Coin struct {
	Code string  `bson:"code" json:"code"`
	Name string  `bson:"name" json:"name"`
	Rate float64 `bson:"rate" json:"rate"`
}

type CoinList []Coin

func (coins CoinList) Len() int           { return len(coins) }
func (coins CoinList) Less(i, j int) bool { return coins[i].Code < coins[j].Code }
func (coins CoinList) Swap(i, j int)      { coins[i], coins[j] = coins[j], coins[i] }
