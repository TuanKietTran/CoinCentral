package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/TuanKietTran/CoinCentral/object"
)

const URL = "https://coin-central-backend.herokuapp.com/"

var state = false

func GetStatus() {
	//TESTING FUNCTION:
	response, err := http.Get(URL + "/status")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		panic("Server may crash!")
	}
	fmt.Println(string(responseData))
}

func GetUser(userID string, platform string) (object.User, bool) {
	//GET USER: FAIL
	api := URL + "users" + "?" + "id=" + userID + "&platform=" + platform
	log.Print("api\n", api, "\n")
	resp, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
		return object.User{}, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return object.User{}, false
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if len(bodyBytes) == 0 {
		return object.User{}, false
	}

	// Convert response body to string
	var user object.User
	json.Unmarshal(bodyBytes, &user)
	return user, true
}

func CreateUser(user object.User) bool {
	//CREATE USER: DONE
	jsonReq, _ := json.Marshal(&user)
	resp, err := http.Post(
		URL+"users",
		"application/json; charset=utf-8",
		bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln(err)
		return false
	}

	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func DeleteUser(userID string, platform string) {
	//DELETE USER: DONE
	api := URL + "users" + "?id=" + userID + "&platform=" + platform
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", api, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Display Results
	fmt.Println("response Status : ", resp.Status)
}

func GetAllCoins() []string {
	// GET ALL COISN: DONE
	api := URL + "coins"
	resp, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	coinArray := strings.Split(bodyString[1:len(bodyString)-1], ",")
	for i := 0; i < len(coinArray); i++ {
		coinArray[i] = coinArray[i][1 : len(coinArray[i])-1]
	}
	return coinArray
}

func GetCoin(coin_name string) (object.Coin, bool) {
	//GET COIN: DONE
	api := URL + "coins/" + coin_name
	fmt.Println(api)
	resp, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Requets error, ", resp.StatusCode)
		return object.Coin{}, false
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if len(bodyBytes) == 0 {
		fmt.Println("Body is empty, Coins may not exists!")
		return object.Coin{}, false
	}

	var coin object.Coin
	err = json.Unmarshal(bodyBytes, &coin)
	if err != nil {
		fmt.Println("Error when API try to convert string to json!")
		return object.Coin{}, false
	}

	return coin, true
}

func GetFollowCoins(userid string, platform string) ([]string, bool) {
	api := URL + "notifications/time?" +
		"id=" + userid +
		"&platform=" + platform +
		"&getCode=true" +
		"&getCode=true" +
		"&getTime=false"

	resp, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Requets error, ", resp.StatusCode)
		return nil, false
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if len(bodyBytes) == 0 {
		fmt.Println("Body is empty, Coins may not exists!")
		return nil, false
	}

	var results = map[string][]string{
		"codeList": []string{},
		"timeList": []string{},
	}
	fmt.Println(">>", string(bodyBytes))
	err = json.Unmarshal(bodyBytes, &results)
	if err != nil {
		fmt.Println("Error when API try to convert string to json!")
		return nil, false
	}
	fmt.Println(results["codeList"])

	return results["codeList"], true
}

func SetFollowCoin(user object.User, coin_code string) bool {
	//TODO: Check format time
	api := URL + "notifications/time?" +
		"id=" + user.Id +
		"&platform=" + user.Platform +
		"&code=" + coin_code

	client := &http.Client{}

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, api, nil)
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for json
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func DeleteFollowCoin(user object.User, coin_code string) bool {
	//TODO: Check format time
	api := URL + "notifications/time?" +
		"id=" + user.Id +
		"&platform=" + user.Platform +
		"&code=" + coin_code

	client := &http.Client{}
	// set the HTTP method, url, and request body
	req, err := http.NewRequest("DELETE", api, nil)
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for json
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func SetTime(user object.User, time string) bool {
	//TODO: Check format time
	api := URL + "notifications/time?" +
		"id=" + user.Id +
		"&platform=" + user.Platform +
		"&time=" + time

	client := &http.Client{}

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, api, nil)
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for json
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func DeleteTime(user object.User, time string) bool {
	//TODO: Check format time
	api := URL + "notifications/time?" +
		"id=" + user.Id +
		"&platform=" + user.Platform +
		"&time=" + time

	client := &http.Client{}
	// set the HTTP method, url, and request body
	req, err := http.NewRequest("DELETE", api, nil)
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for json
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func GetTime(userid string, platform string) ([]string, bool) {
	api := URL + "notifications/time?" +
		"id=" + userid +
		"&platform=" + platform +
		"&getCode=false" +
		"&getTime=true"
	resp, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
		return []string{}, false
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Requets error, ", resp.StatusCode)
		return []string{}, false
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if len(bodyBytes) == 0 {
		fmt.Println("Body is empty, Coins may not exists!")
		return []string{}, false
	}

	var results = map[string][]string{
		"codeList": []string{},
		"timeList": []string{},
	}
	fmt.Println(">>", string(bodyBytes))
	err = json.Unmarshal(bodyBytes, &results)
	if err != nil {
		fmt.Println("Error when API try to convert string to json!")
		return nil, false
	}
	fmt.Println(results["timeList"])
	times := results["timeList"]

	return times, true
}

func SetLimit(user object.User, limit object.Limit) bool {
	//NEED TO CHECK
	api := URL + "notifications/limits" + "?id=" + user.Id + "&platform=" + user.Platform
	jsonReq, err := json.Marshal(&limit)
	fmt.Println("api", api)
	fmt.Println("json", jsonReq)
	if err != nil {
		panic("Limit object cannot parse to json")
	}

	resp, err := http.Post(
		api,
		"application/json; charset=utf-8",
		bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func UpdateLimit(user object.User, limit object.Limit) bool {
	//UPDATE LIMIT : DONE
	api := URL + "notifications/limits" + "?id=" + user.Id + "&platform=" + user.Platform
	jsonReq, err := json.Marshal(&limit)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, api, bytes.NewBuffer(jsonReq))
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func GetLimit(user object.User, coin_code string, isUpper bool) ([]object.Limit, bool) {
	//GET LIMIT DONE
	var mode string
	if isUpper {
		mode = "true"
	} else {
		mode = "false"
	}

	api := URL + "notifications/limits" + "?id=" + user.Id + "&platform=" + user.Platform + "&code=" + coin_code + "&isUpper=" + mode
	resp, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		return nil, false
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if len(bodyBytes) == 0 {
		fmt.Println("Limit does not exists!")
		return nil, false
	}

	var limits []object.Limit
	err = json.Unmarshal(bodyBytes, &limits)
	if err != nil {
		panic(err)
	}
	return limits, true
}

func CreateWebhookRequest(request object.WebhookRequest) bool {
	//Done, require Webhook client to turn on.
	api := URL + "webhook/create"
	jsonReq, err := json.Marshal(&request)
	fmt.Println(string(jsonReq))
	resp, err := http.Post(
		api,
		"application/json;",
		bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

//func main() {
//Check server status:
//getStatus()
// WEBHOOK
// request := WebhookRequest{"http://localhost:8080/", "telegram"}
// success := createWebhookRequest(request)
// if success {
// 	//TODO: Do something
// }
// USER:
//user := User{"123456", "telegram", "Justin Nguyen"}
// success := createUser(user)
// if success {
// 	//TODO: Do something
// }
// deleteUser("123456", "telegram")

// COINS:
// coins := getAllCoins()
// coin, success := getCoin("MANA")
// if success {
// 	//TODO: Do something
// 	fmt.Println(coin)
// }

// LIMIT:
// limit := Limit{"MANA", false, 1.3722288665475049}
// setLimit(user, limit)
// limits, success := getLimit(user, "MANA", true)
// if success == true {
// 	//TODO: Do something.
// 	fmt.Println(limits)
// }

// success := updateLimit(user, limit)
// if success {
// 	//TODO: Do something
// 	fmt.Println("Update success")
// }
//}
