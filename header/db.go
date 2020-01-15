package header

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"strconv"
	"strings"
	"time"
)

var db *bolt.DB
var buckets = []string{"rates"}

func dbInit() {
	var err error
	db, err = bolt.Open("boltDB", 444, nil)
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		for _, name := range buckets {
			_, err = tx.CreateBucketIfNotExists([]byte(name))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	db.Close()
}

type FormattedRate struct {
	Pair      string
	Exchange  string
	BuyPrice  string
	SellPrice string
	Updated   string
}

func (fRate *FormattedRate) FromRate(market CryptoMarket, rate Rate) {
	fRate.Pair = string(rate.CurrPair.First)+"/"+
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

func SaveRate(market CryptoMarket, rate Rate) error {
	var fRate FormattedRate
	fRate.FromRate(market, rate)

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.Bucket([]byte("rates")).CreateBucketIfNotExists([]byte(market.GetName()))
		if err != nil {
			return err
		}
		bytes, err := json.Marshal(fRate)
		if err != nil {
			log.Println(err)
			return err
		}
		return b.Put([]byte(rate.CurrPair.String()), bytes)
	})
	if err != nil {
		log.Println(err)
	}
	return err
}

func LoadRate(market CryptoMarket, pair CurrPair) (Rate, error) {
	var fRate FormattedRate
	err := db.View(func (tx *bolt.Tx) error {
		b := tx.Bucket([]byte("rates")).Bucket([]byte(market.GetName()))
		if b == nil {
			err := errors.New(fmt.Sprintf("no such bucket %v", market.GetName()))
			log.Println(err)
			return err
		}

		bytes := b.Get([]byte(pair.String()))
		if len(bytes) == 0 {
			err := errors.New(fmt.Sprintf("no %v - %v in DB", market.GetName(), pair))
			log.Println(err)
			return err
		}
		if err := json.Unmarshal(bytes, &fRate); err != nil {
			log.Println(err)
			return err
		}
		return nil
	})

	if err != nil {
		log.Println(err)
		return Rate{}, err
	}
	return fRate.ToRate()
}