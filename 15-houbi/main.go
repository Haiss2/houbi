package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"houbi/15-houbi/bot"
	storage "houbi/15-houbi/ram"

	"github.com/nntaoli-project/goex"
	houbi "github.com/nntaoli-project/goex/huobi"
)

var symbols = []string{
	"BTC_USDT",
	"ETH_USDT",
}

const (
	duration                = 10 * time.Second // 30 * time.Second
	snoozeTime              = 5 * time.Minute
	minThreshHoldPercentage = 0.02 // 1 equivalent 1%
)

func StandardDeviation(ps []storage.PriceAPI) (sum, mean, sd float64) {
	length := float64(len(ps))
	if length == 0.0 {
		return
	}
	for _, p := range ps {
		sum += p.Price
	}

	mean = sum / length

	for _, p := range ps {
		sd += math.Pow(p.Price-mean, 2)
	}

	sd = math.Sqrt(sd / length)

	return
}

var mapSnooze = make(map[string]int64)

func removeOldPriceRoutine(ram *storage.RamStorage) {
	ticker := time.NewTicker(duration)
	for {
		<-ticker.C
		expiredTs := time.Now().UTC().UnixMilli() - duration.Milliseconds()
		ram.RemoveExpiredData(expiredTs)

	}
}

func monitorPrice(ram *storage.RamStorage, tele *bot.TelegramBot, sym string) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		ps := ram.GetPricesBySymbol(sym)
		if len(ps) == 0 {
			continue
		}
		_, mean, sd := StandardDeviation(ps)
		lastPrice := ps[len(ps)-1]
		v := (lastPrice.Price - mean) / mean * 100
		if math.Abs(v) > minThreshHoldPercentage {
			t, ok := mapSnooze[sym]
			if !ok || t+snoozeTime.Milliseconds() < time.Now().UnixMilli() {
				stg := fmt.Sprintf("%s \n mean: %.2f,\n standard deviation: %.2f,\n last_price: %.2f,\n volatility:  %.5f%%",
					sym, mean, sd, lastPrice.Price, v)
				tele.Notify(stg, "")
				mapSnooze[sym] = time.Now().UnixMilli()
			}

		}

	}
}

func main() {
	ram := storage.NewRamStorage()

	spotWs := houbi.NewSpotWs()
	spotWs.TickerCallback(func(ticker *goex.Ticker) {
		ram.SavePrice(storage.Price{
			Symbol:    ticker.Pair.String(),
			Price:     ticker.Last,
			Timestamp: int64(ticker.Date),
		})
	})

	for _, sym := range symbols {
		spotWs.SubscribeTicker(goex.NewCurrencyPair2(sym))
	}

	tele, err := bot.NewTelegramBot()
	if err != nil {
		log.Fatal("failed to init tele bot")
	}
	tele.Notify("houbi pricing bot start...", "")

	go removeOldPriceRoutine(ram)
	for _, sym := range symbols {
		go monitorPrice(ram, tele, sym)
	}

	for {
	}
}
