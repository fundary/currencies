currencies
==========

Currencies is a library for consuming OpenExchangeRates.org
It has no external dependencies and has only the features that are
needed for the foundary.com project. You can run the go documentation tool
and get an api overview

Requirements
============

All monetary values are treated as the Money struct from 

    github.com/hailocab/i18n-go/money

Usage
=====

Load the module with

    import "github.com/fundary/currencies"

and set the AppID for your application:

    currencies.AppID("my-app-id")

optionally checking for validation errors that are returned by the function.
Spawn the Updater() goroutine passing an interval to update the rates at (depending on the plan)

    go currencies.Updater(1 * time.Hour)

You can then use the GetRate to get the rate for a given currency,
using the standard three letter symbol (ie: Dollars -> USD, Euro -> EUR), and the
Convert function to convert a given amount from one currency to another

    currencies.Convert(money.New(4200, "USD"), "EUR")

will convert 42 US Dollars in Euro based on the latest rates
