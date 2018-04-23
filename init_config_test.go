package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	found := false
	loadConfig()

	for _, e := range bot.config.Exchanges {
		if e.Name == "LiveCoin" {
			found = true
		}
	}
	assert.Equal(t, found, true)
}
