package header

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
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

type CachedRate struct {
	Rate    Rate
	Updated int64
}
type MuxMap struct {
	MuxMap map[CurrPair]CachedRate
	Mux    sync.Mutex
}

type CryptoMarket interface {
	GetName() string
	GetRate(CurrPair, int64) Rate
}

func GetBody(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := errors.New(fmt.Sprintf("%v %v", url, resp.StatusCode))
		log.Println(err)
		return "", err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	bodyString := string(bodyBytes)

	log.Println(fmt.Sprintf("Got body: %v", url))
	return bodyString, nil
}
