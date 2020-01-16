## CryptoRates

Simple project for getting actual rates from cryptocurrency exchanges.

Supported exchanges |
------------------- |
binance.com |
exmo.com |

Application uses [boltDB](https://github.com/boltdb/bolt) to store rates, install it with:

```sh
$ go get github.com/boltdb/bolt
```

### Getting rates
At the start of the program call `header.Init()`.

Set *currency pairs* which rates you want to track:
```go
var pairs = []header.CurrPair{
	{"BTC", "USDT"},
	{"ADA", "ETH"},
	{"ADA", "BTC"},
	{"DCR", "BTC"},
	{"XTZ", "BTC"},
}
```

Set *exchanges* from where the rates will be tracked:
```go
exchanges = []header.CryptoExchange{&exmo.Exmo{}, &binance.Binance{}}
```

Call `utils.GetRates(pairs, exchanges, recency)` with required *recency* of rates (smaller - more often updates, in seconds):
```go
rates, err := utils.GetRates(pairs, exchanges, 10)
```

The following rates format is used (`RFC3339` for updated):
```json
[{"pair":"ADA/ETH","exchange":"binance","buyPrice":0.00025107,"sellPrice":0.00025154,"updated":"2020-01-17T01:14:39+03:00"},{"pair":"BTC/USDT","exchange":"binance","buyPrice":8705.46,"sellPrice":8705.91,"updated":"2020-01-17T01:14:39+03:00"},{"pair":"BTC/USDT","exchange":"exmo","buyPrice":8686.6,"sellPrice":8725.49999999,"updated":"2020-01-17T01:14:39+03:00"},{"pair":"ADA/ETH","exchange":"exmo","buyPrice":0.00025046,"sellPrice":0.00025197,"updated":"2020-01-17T01:14:40+03:00"}]
```
