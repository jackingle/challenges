package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	bitcoinURL = "https://api.coindesk.com/v1/bpi/currentprice.json"
)

type Bitcoin struct {
	Time struct {
		Updated    string    `json:"updated"`
		UpdatedISO time.Time `json:"updatedISO"`
		Updateduk  string    `json:"updateduk"`
	} `json:"time"`
	Disclaimer string `json:"disclaimer"`
	ChartName  string `json:"chartName"`
	Bpi        struct {
		Usd struct {
			Code        string  `json:"code"`
			Symbol      string  `json:"symbol"`
			Rate        string  `json:"rate"`
			Description string  `json:"description"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"USD"`
		Gbp struct {
			Code        string  `json:"code"`
			Symbol      string  `json:"symbol"`
			Rate        string  `json:"rate"`
			Description string  `json:"description"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"GBP"`
		Eur struct {
			Code        string  `json:"code"`
			Symbol      string  `json:"symbol"`
			Rate        string  `json:"rate"`
			Description string  `json:"description"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"EUR"`
	} `json:"bpi"`
}

// Write a program prints to the console the current bit coin exchange value in US Dollars (GET:: https://api.coindesk.com/v1/bpi/currentprice.json)

func main() {
	if err := getBitcoin(); err != nil {
		log.Fatal(err)
	}
}

func getBitcoin() error {
	client := http.DefaultClient

	resp, err := client.Get(bitcoinURL)
	if err != nil {
		return err
	}
	bitcoinResponse := Bitcoin{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &bitcoinResponse); err != nil {
		return err
	}

	currentValue := fmt.Sprintf("The current current Bitcoin exchange value in US Dollars is 1 Bitcoin for $%v USD.", bitcoinResponse.Bpi.Usd.Rate)
	fmt.Println(currentValue)

	return nil
}
