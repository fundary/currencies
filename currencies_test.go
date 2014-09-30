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
		"EUR": 0.7919,
		"GBP": 0.6166,
		"USD": 1.0,
	},
}

func TestAppIdWrongKey(t *testing.T) {
	err := AppID("someappid")
	if err == nil {
		t.Fatalf("Regexp validation failed")
	}
}

func TestAppIdEmptyKey(t *testing.T) {
	err := AppID("")
	if err == nil {
		t.Fatalf("Accepted empty app_id")
	}
}

func TestGetRate(t *testing.T) {
	current = testExchangeRates
	rate, err := GetRate("EUR")
	if err != nil || rate != 0.7919 {
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

func TestConvertUSDToUSD(t *testing.T) {
	converted, _ := Convert("USD", "USD", 1.0)
	if converted != 1.0 {
		t.Fatalf("Converting USD to USD failed")
	}
}

func TestConvertWrongValues(t *testing.T) {
	c, err := Convert("", "USD", 1.0)
	if c != 0.0 || err == nil {
		t.Fatalf("Converting with empty FROM succeded")
	}
	c, err = Convert("USD", "", 1.0)
	if c != 0.0 || err == nil {
		t.Fatalf("Converting with empty TO succeded")
	}
	c, err = Convert("USD", "EUR", 0.0)
	if c != 0.0 || err == nil {
		t.Fatalf("Converting with amount = 0 succeded")
	}
	c, err = Convert("USD", "EUR", -42.0)
	if c != 0.0 || err == nil {
		t.Fatalf("Converting with amount < 0 succeded")
	}
}

func TestConvertEURToUSD(t *testing.T) {
	converted, _ := Convert("EUR", "USD", 1.0)
	if float64(int(converted*100))/100 != 1.26 {
		t.Fatalf("Converting 1 EUR to USD: %.32f", converted)
	}
}

func TestConvertGBPToUSD(t *testing.T) {
	converted, _ := Convert("GBP", "USD", 1.0)
	if float64(int(converted*100))/100 != 1.62 {
		t.Fatalf("Converting 1 GBP to USD: %f", converted)
	}
}

func TestConvertUSDInEUR(t *testing.T) {
	converted, _ := Convert("USD", "EUR", 1.0)
	if float64(int(converted*100))/100 != 0.79 {
		t.Fatalf("Converting 1 USD to EUR: %f", converted)
	}
}

func TestConvertGBPInEUR(t *testing.T) {
	converted, _ := Convert("GBP", "EUR", 1.0)
	if float64(int(converted*100))/100 != 1.28 {
		t.Fatalf("Converting 1 GBP to EUR: %f", converted)
	}
}
