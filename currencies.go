// Package currencies provides a simple interface to the http://openexchangerates.org api
package currencies

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
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

// Rate is used to wrap a single exchange rate value
type Rate struct {
	Name  string
	Value int64
	When  time.Time
}

// Get Latest returns the latest rates.
func (c *OpenExchangeClient) GetLatest() (rates []Rate, err error) {
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

	t := time.Unix(ex.Timestamp, 0)
	for name, cur := range ex.Rates {
		rates = append(rates, Rate{Name: name, Value: int64(cur * 100), When: t})
	}
	return
}
