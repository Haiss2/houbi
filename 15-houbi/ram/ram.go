package storage

import "sync"

type ramDB map[string][]PriceAPI

type Price struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

type PriceAPI struct {
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

func (p Price) ToAPI() PriceAPI {
	return PriceAPI{p.Price, p.Timestamp}
}

func (p PriceAPI) ToPrice(symbol string) Price {
	return Price{symbol, p.Price, p.Timestamp}
}

type RamStorage struct {
	mu sync.RWMutex
	db ramDB
}

func NewRamStorage() *RamStorage {
	db := make(ramDB, 0)
	return &RamStorage{
		mu: sync.RWMutex{},
		db: db,
	}
}

func (r *RamStorage) SavePrice(p Price) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.db[p.Symbol] = append(r.db[p.Symbol], p.ToAPI())
}

func (r *RamStorage) SaveBatchPrice(prices []Price) {
	for _, p := range prices {
		r.SavePrice(p)
	}
}

func (r *RamStorage) GetPricesBySymbol(symbol string) []PriceAPI {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.db[symbol]
}

func (r *RamStorage) RemoveExpiredData(timestamp int64) []Price {
	r.mu.Lock()
	defer r.mu.Unlock()
	removed := make([]Price, 0)
	for symbol, prices := range r.db {
		start := 0
		for i, p := range prices {
			if p.Timestamp > timestamp {
				start = i
				break
			} else {
				removed = append(removed, p.ToPrice(symbol))
				// if p is the last of slice and p is expired
				if i == len(prices)-1 {
					start = i + 1
				}
			}

		}
		copy(r.db[symbol], r.db[symbol][start:])
		r.db[symbol] = r.db[symbol][:len(r.db[symbol])-start]
	}
	return removed
}
