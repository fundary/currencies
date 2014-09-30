// Package currencies provides a simple interface to the http://openexchangerates.org api
// It can be used to keep a local copy of the latest state of the api.
//
// Usage:
// First, you need to import the package as usual and then set the AppId by calling the
// AppId function:
//				AppId("s0m34pp1d")
// Then you can spawn the Updater as a goroutine that will poll the api every
// "interval" of time and update the local state
package currencies

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

const (
	apiBase    = "http://openexchangerates.org/api"
	latest     = apiBase + "/latest.json?app_id="
	currencies = apiBase + "/currencies?app_id="
)

var (
	appId   string
	client  *http.Client
	current *ExchangeRate
)

func init() {
	client = &http.Client{}
}

// Type Exchangerate is used to marshal the current rates from the json api
// at https://openexchangerates.org/documentation
type ExchangeRate struct {
	Disclaimer string             `json:"disclaimer"`
	License    string             `json:"license"`
	Timestamp  int64              `json:"timestamp"`
	Base       string             `json:"base"`
	Rates      map[string]float64 `json:"rates"`
}

// Type Rate is used to wrap a single exchange rate value
type Rate struct {
	Name  string
	Value float64
}

// Function AppId sets the local app_id key used to access the api. It validates
// the api key via regexp, so no test of functioning is done and a wrong api key
// can lead to runtime errors
func AppId(id string) (err error) {
	pattern := "[a-zA-Z0-9]{32}"
	ok, err := regexp.MatchString(pattern, id)
	if err != nil {
		return fmt.Errorf("Error validating regexp: %s", err)
	}
	if !ok {
		return fmt.Errorf("Invalid app_id provided")
	}
	appId = id
	return nil
}

func doRequest(url string) (data []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return
}

func setNewRates(ex *ExchangeRate) {
	current = ex
}

// Function Latest returns an ExchangeRates object loaded with the latest
// rates fetched from the api
func Latest() (ex *ExchangeRate, err error) {
	ex = new(ExchangeRate)
	log.Println("Updating exchange rates")
	url := latest + appId
	fmt.Println(url)
	data, err := doRequest(url)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &ex)
	if err != nil {
		return nil, err
	}
	return
}

// Function Updater is launched as a goroutine and updates the current
// vision of the exchange rates of the library from the api
func Updater(interval time.Duration) {
	log.Printf("Fetching new currencies every %s\n", interval)
	ticker := time.NewTicker(interval)
	select {
	case <-ticker.C:
		ex, err := Latest()
		if err != nil {
			log.Printf("Could not update rates: %s", err)
		} else {
			setNewRates(ex)
		}
	}
}

// Function GetRate queries the current state of the view of the exchange rates
// and returns the rate for the given currency if available or an error if the
// selected currency does not exist
func GetRate(cur string) (rate float64, err error) {
	rate, ok := current.Rates[cur]
	if !ok {
		return rate, fmt.Errorf("Currency %s not available", cur)
	}
	return rate, nil
}

// Function GetRates returns a buffered channel that holds at most the current
// number of exchange rates. A slow reading client will receive all the rates
// as they were at the moment the function was called (a copy of the struct is passed to
// the inner goroutine)
func GetRates() <-chan (Rate) {
	rates := make(chan (Rate), len(current.Rates))
	go func(cur ExchangeRate, c chan (Rate)) {
		defer close(c)
		for k, v := range cur.Rates {
			c <- Rate{k, v}
		}
	}(*current, rates)
	return rates
}
