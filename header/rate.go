package header

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type Rate struct {
	CurrPair  CurrPair
	BuyPrice  float64
	SellPrice float64
	Updated   int64
}

func (rate Rate) String() string {
	return fmt.Sprintf("%v buy: %v; sell: %v; updated: %v\n",
		rate.CurrPair, rate.BuyPrice, rate.SellPrice, time.Unix(rate.Updated, 0))
}

type FormattedRate struct {
	Pair      string
	Exchange  string
	BuyPrice  float64
	SellPrice float64
	Updated   string
}

func (fRate *FormattedRate) FromRate(market CryptoMarket, rate Rate) {
	fRate.Pair = string(rate.CurrPair.First) + "/" +
	    string(rate.CurrPair.Second)
	fRate.Exchange = market.GetName()
	fRate.BuyPrice = rate.BuyPrice
	fRate.SellPrice = rate.SellPrice
	fRate.Updated = time.Unix(rate.Updated, 0).Format(time.RFC3339)
}

func (fRate *FormattedRate) ToRate() (Rate, error) {
	var rate Rate
	var err error

	s := strings.Split(fRate.Pair, "/")
	if len(s) != 2 {
		err = errors.New(fmt.Sprintf("couldn't convert to Rate: %v", fRate))
		log.Println(err)
		return Rate{}, err
	}
	rate.CurrPair = CurrPair{s[0], s[1]}
	rate.BuyPrice = fRate.BuyPrice
	rate.SellPrice = fRate.SellPrice

	updatedTime, err := time.Parse(time.RFC3339, fRate.Updated)
	if err != nil {
		log.Printf("couldn't parse updated %v", fRate.Updated)
		return Rate{}, err
	}
	rate.Updated = updatedTime.Unix()
	return rate, nil
}
