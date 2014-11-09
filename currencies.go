// Package currencies provides a simple interface to the http://openexchangerates.org api
package currencies

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	apiBase    = "https://openexchangerates.org/api"
	latest     = apiBase + "/latest.json?app_id="
	currencies = apiBase + "/currencies?app_id="
)

type OpenExchangeClient struct {
	AppID  string
	Client *http.Client
}

// ExchangeRate is used to marshal the current rates from the json api
// at https://openexchangerates.org/documentation
type exchangeRate struct {
	Disclaimer string             `json:"disclaimer"`
	License    string             `json:"license"`
	Timestamp  int64              `json:"timestamp"`
	Base       string             `json:"base"`
	Rates      map[string]float64 `json:"rates"`
}

// Get Latest returns the latest rates.
func (c *OpenExchangeClient) GetLatest() (rates map[string]int64, err error) {
	if c.AppID == "" {
		return nil, errors.New("No app Id")
	}
	if c.Client == nil {
		c.Client = &http.Client{}
	}

	resp, err := c.Client.Get(latest + c.AppID)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()

	ex := new(exchangeRate)
	dec := json.NewDecoder(resp.Body)

	if err := dec.Decode(&ex); err == io.EOF {

	} else if err != nil {
		return nil, err
	}
	rates = make(map[string]int64)
	for name, cur := range ex.Rates {
		rates[name] = int64(cur * 100)
	}
	return
}
