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

var getRates = make(map[string]func(CurrPair) (string, error))
var getUrl = make(map[string]func(CurrPair) string)

func Init() {
	getUrl["exmoUrl"] = func (currPair CurrPair) string {
		return fmt.Sprintf("https://api.exmo.com/v1/trades/?pair=%v", currPair)
	}

	getRates["exmoRates"] = func(currPair CurrPair) (string, error) {
		url, err := GetUrl("exmoUrl", currPair)
		if err != nil {
			return "", err
		}

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
				return "", err
			}
			bodyString := string(bodyBytes)

			log.Println("Got rates ", url)
			return bodyString, nil
		} else {
			err := errors.New(fmt.Sprintf("%v %v", url, resp.StatusCode))
			log.Fatal(err)
			return "", err
		}

		return "", nil
	}
}

func GetUrl(site string, currPair CurrPair) (string, error) {
	f, ok := getUrl[site]
	if !ok {
		err := errors.New(fmt.Sprintf("%v isn't supported", site))
		return "", err
	}
	url := f(currPair)
	return url, nil
}
func GetRates(site string, currPair CurrPair) (string, error) {
	f, ok := getRates[site]
	if !ok {
		err := errors.New(fmt.Sprintf("%v isn't supported", site))
		return "", err
	}

	res, err := f(currPair)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	fmt.Println(res)
	if res == "{}" {
		err := errors.New(fmt.Sprintf("%v doesn't support %v", site, currPair))
		log.Fatal(err)
		return "", nil
	}
	return res, nil
}

func main() {
	Init()

	res, err := GetRates("exmoRates", CurrPair{"BTC", "USD"})
	fmt.Println(err)
	fmt.Println(res)
}
