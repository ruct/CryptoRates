package header

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
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

func SaveRate(exchange CryptoExchange, rate Rate) error {
	var fRate FormattedRate
	fRate.FromRate(exchange, rate)

	err := db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.Bucket([]byte("rates")).CreateBucketIfNotExists([]byte(exchange.GetName()))
		if err != nil {
			return err
		}
		bytes, err := json.Marshal(fRate)
		if err != nil {
			log.Println(err)
			return err
		}
		return b.Put([]byte(rate.Pair.String()), bytes)
	})
	if err != nil {
		log.Println(err)
	}
	return err
}

func LoadRate(market CryptoExchange, pair CurrPair) (Rate, error) {
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