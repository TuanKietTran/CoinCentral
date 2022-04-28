package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Coins struct {
		NumOfSupportingCoins int `yaml:"numOfSupportingCoins"`
		NumOfFetchCoin       int `yaml:"numOfFetchCoin"`
		TimeBetweenFetch     int `yaml:"timeBetweenFetch"`
	} `yaml:"coins"`
}

func ReadConfigFile(configFileName string) *Config {
	configFile, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Panicf("Can't read config file, %v", err)
	}

	var config Config
	if err = yaml.Unmarshal(configFile, &config); err != nil {
		log.Panicf("Can't parse config file, %v", err)
	}

	return &config
}
