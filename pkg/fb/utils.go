package fb

import "log"

func setBounds(userData UserMapIDKey, recipientID string, coinName string, lower int, upper int) error {
	if temp, existed := userData[recipientID]; existed {
		temp.State = "pend"

		if entry, ok := userData[recipientID].Data.FollowedCoinList[coinName]; ok {
			entry.Lowerbound = lower
			entry.Upperbound = upper
			entry.Name = coinName
			userData[recipientID].Data.FollowedCoinList[coinName] = entry
		} else {
			log.Println("set bound function")
		}
	}
	// if entry, ok := userData.Data.FollowedCoinList[coinName]; ok {
	// 	entry.Lowerbound = lower
	// 	entry.Upperbound = upper
	// 	entry.Name = coinName
	// 	userData.Data.FollowedCoinList[coinName] = entry
	// } else {
	// 	log.Println("not exist so create new coinData named ", coinName)
	// 	userData.Data.FollowedCoinList = make(CoinMapNameKey)
	// 	userData.Data.FollowedCoinList[coinName] = CoinData{
	// 		Name:       coinName,
	// 		Lowerbound: lower,
	// 		Upperbound: upper,
	// 	}

	// 	log.Println(userData.Data.FollowedCoinList)
	// }

	return nil
}
