package header

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
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
	Updated int64
}

func (rate Rate) String() string {
	return fmt.Sprintf("%v buy: %v; sell: %v; updated: %v\n", rate.CurrPair, rate.BuyPrice, rate.SellPrice, time.Unix(rate.Updated, 0))
}

type CryptoMarket interface {
	GetName() string
	GetRate(CurrPair, int64) (Rate, error)
	GetTradesUrl(CurrPair) string
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

	return bodyString, nil
}

func GetJson(url string) (map[string]interface{}, error) {
	body, err := GetBody(url)
	if err != nil {
		return nil, err
	}

	bytes := []byte(body)
	var fullData map[string]interface{}
	if err := json.Unmarshal(bytes, &fullData); err != nil {
		log.Println(err)
		return nil, err
	}
	return fullData, nil
}

func Init() {
	runtime.GOMAXPROCS(8)

	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)
	log.Println("log inited")
}


func DefaultGetRate(market CryptoMarket, currPair CurrPair, recency int64,
	getCachedRate func() (Rate, bool), renew func() error) (Rate, error) {

	cachedRate, ok := getCachedRate()
	if !ok {
		log.Println(fmt.Sprintf("%v: no such %v", market.GetName(), currPair))
	}
	if ok && time.Now().Unix()-cachedRate.Updated > recency {
		log.Println(fmt.Sprintf("%v: need to update %v, last update was: %v", market.GetName(), currPair, time.Now().Unix()-cachedRate.Updated))
	}

	if !ok || time.Now().Unix()-cachedRate.Updated > recency {
		was := cachedRate.Updated

		err := renew()
		if err != nil {
			log.Println(err)
			return Rate{}, err
		}

		cachedRate, ok = getCachedRate()
		if !ok {
			err := errors.New(fmt.Sprintf("%v: incorrect pair %v", market.GetName(), currPair))
			log.Println(err)
			return Rate{}, err
		}

		became := cachedRate.Updated
		log.Println(fmt.Sprintf("%v: wanted %v, was: %v, became: %v", market.GetName(), currPair, was, became))
	}
	return cachedRate, nil
}

func DefaultRenew(market CryptoMarket, currPair CurrPair,
	processJson func (map[string]interface{}) error) error {

	log.Println(fmt.Sprintf("%v: updating %v", market.GetName(), time.Now()))
	fmt.Println(fmt.Sprintf("%v: updating %v", market.GetName(), time.Now()))
	fullData, err := GetJson(market.GetTradesUrl(currPair))
	if err != nil {
		return err
	}
	return processJson(fullData)
}