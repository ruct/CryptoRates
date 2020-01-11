package exmo

import (
	"../header"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Exmo struct {
	cachedRates header.MuxMap
}

func (*Exmo) GetName() string {
	return "exmo"
}

func (exmo *Exmo) GetRate(currPair header.CurrPair, recency int64) (header.Rate, error) {
	exmo.cachedRates.Mux.Lock()
	defer exmo.cachedRates.Mux.Unlock()

	cachedRate, ok := exmo.cachedRates.MuxMap[currPair]
	if !ok || time.Now().Unix()-cachedRate.Updated > recency {
		exmo.renew()
		cachedRate, ok = exmo.cachedRates.MuxMap[currPair]
		if !ok {
			err := errors.New(fmt.Sprintf("exmo: no such %v", currPair))
			return header.Rate{}, err
		}
	}
	return cachedRate.Rate, nil
}

func (exmo *Exmo) renew() error {
	body, err := header.GetBody("https://api.exmo.com/v1/ticker/")
	if err != nil {
		return err
	}

	bytes := []byte(body)
	var fullData map[string]interface{}
	if err := json.Unmarshal(bytes, &fullData); err != nil {
		log.Println(err)
		return err
	}

	var data = make(map[header.CurrPair]header.CachedRate)
	for key, value := range fullData {
		s := strings.Split(key, "_")
		if len(s) != 2 {
			log.Println(fmt.Sprintf("exmo: couldn't convert currency-pair %v", key))
			continue
		}
		currPair := header.CurrPair{s[0], s[1]}

		var cmap = value.(map[string]interface{})

		buyPrice, ok1 := cmap["buy_price"].(string)
		sellPrice, ok2 := cmap["sell_price"].(string)
		updated, ok3 := cmap["updated"].(float64)
		if !(ok1 && ok2 && ok3) {
			switch {
			case !ok1:
				log.Println(fmt.Sprint("exmo: couldn't parse %v: buyPrice %v"), cmap["buy_price"])
			case !ok2:
				log.Println(fmt.Sprint("exmo: couldn't parse sellPrice %v"), cmap["sell_price"])
			case !ok3:
				log.Println(fmt.Sprint("exmo: couldn't parse updated %v"), cmap["updated"])
			}
			continue
		} else {
			buyPrice, err1 := strconv.ParseFloat(buyPrice, 64)
			sellPrice, err2 := strconv.ParseFloat(sellPrice, 64)
			updated := int64(updated)

			if err1 != nil {
				log.Println(err1)
				return err1
			}
			if err2 != nil {
				log.Println(err2)
				return err2
			}

			data[currPair] = header.CachedRate{
				header.Rate{
					currPair,
					buyPrice,
					sellPrice,
				},
				updated}
		}
	}

	exmo.cachedRates.MuxMap = data
	return nil
}
