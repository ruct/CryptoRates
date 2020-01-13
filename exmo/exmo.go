package exmo

import (
	"../header"
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

	return header.DefaultGetRate(exmo, currPair, recency,
		func() (rate header.Rate, ok bool) {
			rate, ok = exmo.cachedRates.MuxMap[currPair]
			return rate, ok
		},
		func() error {
			return exmo.renew()
		})
}

func (exmo *Exmo) GetTradesUrl(header.CurrPair) string {
	return "https://api.exmo.com/v1/ticker/"
}

func (exmo *Exmo) processJson(jsonData map[string]interface{}) error {
	if exmo.cachedRates.MuxMap == nil {
		exmo.cachedRates.MuxMap = make(map[header.CurrPair]header.Rate)
	}

	var data = make(map[header.CurrPair]header.Rate)
	for key, value := range jsonData {
		s := strings.Split(key, "_")
		if len(s) != 2 {
			log.Println(fmt.Sprintf("%v: couldn't convert currency-pair %v", exmo.GetName(), key))
			continue
		}
		currPair := header.CurrPair{s[0], s[1]}

		var cmap = value.(map[string]interface{})

		buyPrice, ok1 := cmap["buy_price"].(string)
		sellPrice, ok2 := cmap["sell_price"].(string)
		if !(ok1 && ok2) {
			switch {
			case !ok1:
				log.Println(fmt.Sprintf("%v: couldn't parse buyPrice %v", exmo.GetName(), cmap["buy_price"]))
			case !ok2:
				log.Println(fmt.Sprintf("%v: couldn't parse sellPrice %v", exmo.GetName(), cmap["sell_price"]))
			}
			continue
		} else {
			buyPrice, err1 := strconv.ParseFloat(buyPrice, 64)
			sellPrice, err2 := strconv.ParseFloat(sellPrice, 64)

			if err1 != nil {
				log.Println(err1)
				return err1
			}
			if err2 != nil {
				log.Println(err2)
				return err2
			}

			data[currPair] = header.Rate{
				currPair,
				buyPrice,
				sellPrice,
				time.Now().Unix(),
			}
		}
	}

	exmo.cachedRates.MuxMap = data
	return nil
}

func (exmo *Exmo) renew() error {
	return header.DefaultRenew(exmo, header.CurrPair{},
		func(jsonData map[string]interface{}) error {
			return exmo.processJson(jsonData)
		})
}
