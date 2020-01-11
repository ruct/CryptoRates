package exmo

import (
	"../header"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

type cachedRate struct {
	rate    header.Rate
	recency int64
}
type muxMap struct {
	muxMap map[header.CurrPair]cachedRate
	mux    sync.Mutex
}
type Exmo struct {
	cachedRates muxMap
}

func (*Exmo) GetName() string {
	return "exmo"
}

func (exmo *Exmo) GetRate(currPair header.CurrPair) (header.Rate, error) {
	return header.Rate{}, nil
}

func (exmo *Exmo) Renew() error {
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

	var data = make(map[header.CurrPair]cachedRate)
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

			data[currPair] = cachedRate{
				header.Rate{
					currPair,
					buyPrice,
					sellPrice,
				},
				updated}
			fmt.Println(fmt.Sprintf("%v: %v %v %v", currPair, buyPrice, sellPrice, updated))
		}
	}

	exmo.cachedRates.mux.Lock()
	exmo.cachedRates.muxMap = data
	exmo.cachedRates.mux.Unlock()

	return nil
}
