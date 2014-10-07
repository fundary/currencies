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

	"github.com/hailocab/i18n-go/money"
)

const (
	apiBase    = "http://openexchangerates.org/api"
	latest     = apiBase + "/latest.json?app_id="
	currencies = apiBase + "/currencies?app_id="
)

var (
	appID   string = ""
	client  *http.Client
	current *ExchangeRate
)

func init() {
	client = &http.Client{}
}

// ExchangeRate is used to marshal the current rates from the json api
// at https://openexchangerates.org/documentation
type ExchangeRate struct {
	Disclaimer string             `json:"disclaimer"`
	License    string             `json:"license"`
	Timestamp  int64              `json:"timestamp"`
	Base       string             `json:"base"`
	Rates      map[string]float64 `json:"rates"`
}

// Rate is used to wrap a single exchange rate value
type Rate struct {
	Name  string
	Value float64
}

// AppID sets the local app_id key used to access the api. It validates
// the api key via regexp, so no test of functioning is done and a wrong api key
// can lead to runtime errors
func AppID(id string) (err error) {
	pattern := "[a-zA-Z0-9]{32}"
	ok, err := regexp.MatchString(pattern, id)
	if err != nil {
		return fmt.Errorf("Error validating regexp: %s", err)
	}
	if !ok {
		return fmt.Errorf("Invalid app_id provided")
	}
	appID = id
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

func validCurrency(cur string) bool {
	_, valid := current.Rates[cur]
	return valid
}

func getLatest() (ex *ExchangeRate, err error) {
	if appID == "" {
		panic("Required AppID() not configured!")
	}
	ex = new(ExchangeRate)
	log.Println("Updating exchange rates")
	url := latest + appID
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

// Update is updates the current vision of the exchange rates
// of the library from the net
func Update() (err error) {
	log.Println("Fetching new currencies")
	ex, err := getLatest()
	if err != nil {
		log.Printf("Could not update rates: %s", err)
		return err
	}
	setNewRates(ex)
	return nil
}

// GetRate queries the current state of the view of the exchange rates
// and returns the rate for the given currency if available or an error if the
// selected currency does not exist
func GetRate(cur string) (rate float64, err error) {
	rate, ok := current.Rates[cur]
	if !ok {
		return rate, fmt.Errorf("Currency %s not available", cur)
	}
	return rate, nil
}

// GetRates returns a buffered channel that holds at most the current
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

// Convert converts from one currency to another, returning a converted value or an
// error if either curency does not exist
func Convert(amount *money.Money, to string) (converted *money.Money, err error) {
	if !validCurrency(amount.C) {
		return nil, fmt.Errorf("Currency %s does not exist or is not available", amount.C)
	}
	if !validCurrency(to) {
		return nil, fmt.Errorf("Currency %s does not exist or is not available", to)
	}
	fromRate, _ := current.Rates[amount.C]
	toRate, _ := current.Rates[to]
	fromToRate := toRate * (1 / fromRate)
	return money.New(amount.Mulf(fromToRate).Value(), to), nil
}
