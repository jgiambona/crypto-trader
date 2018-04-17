package exchange

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// SendPayload creates a request to be sent to endpoint.
func (o *Base) SendPayload(method, path string, headers map[string]string,
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

	err = json.Unmarshal(contents, &result)
	return err
}
