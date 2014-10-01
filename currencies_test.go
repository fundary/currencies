package currencies

import (
	"log"
	"testing"
	"time"

	"github.com/hailocab/i18n-go/money"
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
	converted, _ := Convert(money.New(100, "USD"), "USD")
	if converted.Get() != 1.0 || converted.C != "USD" {
		t.Fatalf("Converting USD to USD failed: %f", converted.Get())
	}
}

func TestConvertWrongValues(t *testing.T) {
	c, err := Convert(money.New(100, "USD"), "")
	if c != nil || err == nil {
		t.Fatalf("Converting with empty TO succeded")
	}
}

func TestConvertEURToUSD(t *testing.T) {
	converted, _ := Convert(money.New(100, "EUR"), "USD")
	if converted.Get() != 1.26 || converted.C != "USD" {
		t.Fatalf("Converting 1 EUR to USD: %.32f", converted)
	}
}

func TestConvertGBPToUSD(t *testing.T) {
	converted, _ := Convert(money.New(100, "GBP"), "USD")
	if converted.Get() != 1.62 || converted.C != "USD" {
		t.Fatalf("Converting 1 GBP to USD: %f", converted)
	}
}

func TestConvertUSDInEUR(t *testing.T) {
	converted, _ := Convert(money.New(100, "USD"), "EUR")
	if converted.Get() != 0.79 || converted.C != "EUR" {
		t.Fatalf("Converting 1 USD to EUR: %f", converted)
	}
}

func TestConvertGBPInEUR(t *testing.T) {
	converted, _ := Convert(money.New(100, "GBP"), "EUR")
	if converted.Get() != 1.28 || converted.C != "EUR" {
		t.Fatalf("Converting 1 GBP to EUR: %f", converted)
	}
}
