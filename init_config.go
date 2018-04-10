package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type (
	// Wallet TODO
	Wallet struct {
		Address     string
		Type        string
		Balance     float64
		Description string
	}

	// Exchange TODO
	Exchange struct {
		Name      string
		RateLimit int64
		URL       struct {
			WebLocation string
			APILocation string
		}
		APIKey    string
		APISecret string
	}

	// BotConfig TODO
	BotConfig struct {
		Portfolio struct {
			Addresses []Wallet
		}
		Exchanges []Exchange
	}
)

func loadConfig() {
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(data), &bot.config)
	if err != nil {
		panic(err)
	}
}
