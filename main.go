package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CurrPair struct {
	first  string
	second string
}

func (currPair CurrPair) String() string {
	return fmt.Sprintf("%v_%v", currPair.first, currPair.second)
}

type CryptoMarket struct {
	name         string
	getTradesUrl func(CurrPair) string
}

func (market *CryptoMarket) String() string {
	return market.name
}
func (market *CryptoMarket) getTrades(currPair CurrPair) (string, error) {
	if market.getTradesUrl == nil {
		err := errors.New(fmt.Sprintf("%v has no getTradesUrl()", market.name))
		log.Println(err)
		return "", err
	}
	return GetBody(market.getTradesUrl(currPair))
}

var markets []CryptoMarket
var marketByName = make(map[string]*CryptoMarket)

func Init() {
	markets = []CryptoMarket{
		{name: "exmo"},
		{name: "binance"},
	}
	for i := range markets {
		marketByName[markets[i].name] = &markets[i]
	}

	marketByName["exmo"].getTradesUrl = func(currPair CurrPair) string {
		return fmt.Sprintf("https://api.exmo.com/v1/trades/?pair=%v", currPair)
	}
}

func GetBody(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := errors.New(fmt.Sprintf("%v %v", url, resp.StatusCode))
		log.Fatal(err)
		return "", err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	bodyString := string(bodyBytes)

	log.Println("Got body: %v", url)
	return bodyString, nil
}

func main() {
	Init()
	res, err := marketByName["exmo"].getTrades(CurrPair{"BTC", "USTD"})
	fmt.Println(err)
	fmt.Println(res)
}
