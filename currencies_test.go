package currencies_test

import (
	"appengine"
	"appengine/urlfetch"
)

// Create an OpenExchange Client
func ExampleOpenExchangeClient() {
	oe := OpenExchangeClient{
		AppID: "{APP_ID}",
	}
}

// If you are using GAE or have an existing http.Client
// You can use it's context to fetch the rates like so
func ExampleOpenExchangeClient_gae() {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	oe := OpenExchangeClient{
		AppID:  "{APP_ID}",
		Client: client,
	}
}

// Get the latest rates
func ExampleOpenExchangeClient_GetLatest() {
	oe := OpenExchangeClient{
		AppID: "{APP_ID}",
	}
	rates, err := oe.GetLatest()
}

// Get the latest rates
func ExampleOpenExchangeClient_GetLatest_gae() {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	oe := OpenExchangeClient{
		AppID:  "{APP_ID}",
		Client: client,
	}
	rates, err := oe.GetLatest()
}
