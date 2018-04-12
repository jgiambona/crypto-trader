package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type (
	// Wallet TODO
	Wallet struct {
		Address string `yaml:"address"`
		Type    string `yaml:"type"`
	}

	// Exchange TODO
	Exchange struct {
		Name      string `yaml:"name"`
		RateLimit int64  `yaml:"rateLimit"`
		URL       struct {
			WebLocation string `yaml:"web"`
			APILocation string `yaml:"api"`
		} `yaml:"urls"`
		APIKey    string `yaml:"apiKey"`
		APISecret string `yaml:"apiSecret"`
	}

	// BotConfig TODO
	BotConfig struct {
		Portfolio struct {
			Addresses []Wallet
		} `yaml:"portfolio"`
		Exchanges []Exchange `yaml:"exchanges"`
	}
)

func loadConfig() {
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	err = yaml.Unmarshal([]byte(data), &bot.config)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
