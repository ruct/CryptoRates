package header

import (
	"errors"
	"fmt"
	"log"
	"strconv"
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
	BuyPrice  string
	SellPrice string
	Updated   string
}

func (fRate *FormattedRate) FromRate(market CryptoMarket, rate Rate) {
	fRate.Pair = string(rate.CurrPair.First) + "/" +
	    string(rate.CurrPair.Second)
	fRate.Exchange = market.GetName()
	fRate.BuyPrice = strconv.FormatFloat(rate.BuyPrice, 'f', 10, 64)
	fRate.SellPrice = strconv.FormatFloat(rate.SellPrice, 'f', 10, 64)
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
	rate.BuyPrice, err = strconv.ParseFloat(fRate.BuyPrice, 64)
	if err != nil {
		log.Println(fmt.Sprintf("couldn't parse buyPrice %v", fRate.BuyPrice))
		return Rate{}, err
	}

	rate.SellPrice, err = strconv.ParseFloat(fRate.SellPrice, 64)
	if err != nil {
		log.Println(fmt.Sprintf("couldn't parse sellPrice %v", fRate.SellPrice))
		return Rate{}, err
	}

	updatedTime, err := time.Parse(time.RFC3339, fRate.Updated)
	if err != nil {
		log.Println(fmt.Sprintf("couldn't parse updated %v", fRate.Updated))
		return Rate{}, err
	}
	rate.Updated = updatedTime.Unix()
	return rate, nil
}
