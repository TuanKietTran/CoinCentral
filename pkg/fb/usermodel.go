package fb

type UserMapIDKey = map[string]UserData
type CoinMapNameKey = map[string]CoinData

type (
	UserData struct {
		State string //full, pend
		Data  RequestData
	}

	RequestData struct {
		FollowedCoinList CoinMapNameKey
		TimeInterval     int //mins
	}

	CoinData struct {
		Name       string
		Upperbound int
		Lowerbound int
	}
)
