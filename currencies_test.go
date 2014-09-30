package currencies

import (
	"log"
	"testing"
	"time"
)

var testExchangeRates = &ExchangeRate{
	Disclaimer: "Disclaimer",
	License:    "License",
	Timestamp:  time.Now().Unix(),
	Base:       "USD",
	Rates: map[string]float64{
		"USD": 1.0,
		"EUR": 0.7,
		"GBP": 1.2,
	},
}

func TestAppIdWrongKey(t *testing.T) {
	err := AppId("someappid")
	if err == nil {
		t.Fatalf("Regexp validation failed")
	}
}

func TestAppIdEmptyKey(t *testing.T) {
	err := AppId("")
	if err == nil {
		t.Fatalf("Accepted empty app_id")
	}
}

func TestGetRate(t *testing.T) {
	current = testExchangeRates
	rate, err := GetRate("USD")
	if err != nil || rate != 1.0 {
		t.Fatalf("Error retrieving expected value")
	}
	rate, err = GetRate("NONEXISTENT")
	if err == nil {
		t.Fatalf("Allowed fetch of non existing currency")
	}
}

func TestGetRates(t *testing.T) {
	for v := range GetRates() {
		_, ok := current.Rates[v.Name]
		if !ok {
			log.Fatalf("Did not retrieve rates for all currencies")
		}
	}
}
