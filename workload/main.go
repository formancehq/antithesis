package main

import (
	"net/http"

	"github.com/antithesishq/antithesis-sdk-go/lifecycle"
) 

type Details map[string]any


func main() {
	_, err := http.Get("http://ledger:8080/_info")
	if err != nil {
		panic(err)
	} else {
		// signals that the system is up and running
		lifecycle.SetupComplete(Details{"Ledger": "Available"})
	}

	// TODO: use SDK randomness when generating transactions

	// TODO: validate and use SDK assertions
}