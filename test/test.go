package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	minimumOrder = "0.001"
	apiBase      = "https://api.livecoin.net"
	apiKey       = "vGMWuD7nw6WDQsRSA3SxW4mnqe4CnG4x"
	apiSecret    = "dVvBd2HCCmEsJCUHSYCxUgutskkUJEbe"
)

type (
	OrderResponse struct {
		Success bool  `json:"success"`
		Added   bool  `json:"added"`
		OrderID int64 `json:"orderId"`
	}
)

func main() {
	path := fmt.Sprintf("%s%s", apiBase, "/exchange/buylimit")

	construct := url.Values{}
	construct.Add("currencyPair", "BTC/USD")
	construct.Add("price", "60")
	construct.Add("quantity", minimumOrder)
	message := construct.Encode()

	headers := map[string]string{
		"API-Key":        apiKey,
		"Sign":           createSignature(message, apiSecret),
		"Content-Type":   "application/x-www-form-urlencoded",
		"Content-Length": strconv.Itoa(len(message)),
	}

	data := OrderResponse{}
	err := sendPayload("POST", path, headers, strings.NewReader(message), &data)
	if err != nil {
		log.Panic(err.Error())
	}
	log.Printf("Hello World")
}

func sendPayload(method, path string, headers map[string]string,
	body io.Reader, result interface{}) error {
	method = strings.ToUpper(method)

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	log.Println(res.Status)
	log.Println(res.Header)
	log.Println(string(contents))
	err = json.Unmarshal(contents, &result)
	return err
}

func createSignature(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	d := hex.EncodeToString(h.Sum(nil))
	return strings.ToUpper(d)
}
