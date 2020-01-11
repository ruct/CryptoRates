package header

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CurrPair struct {
	First  string
	Second string
}

func (currPair CurrPair) String() string {
	return fmt.Sprintf("%v_%v", currPair.First, currPair.Second)
}

type Rate struct {
	CurrPair  CurrPair
	BuyPrice  float64
	SellPrice float64
}

func (rate Rate) String() string {
	return fmt.Sprintf("%v buy: %v; sell: %v\n", rate.CurrPair, rate.BuyPrice, rate.SellPrice)
}

type CryptoMarket interface {
	GetName() string
	GetRate(CurrPair, int32) Rate
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
