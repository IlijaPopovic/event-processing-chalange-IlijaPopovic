package currency

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"
)

type Converter struct {
	cache     map[string]float64
	cacheLock sync.RWMutex
	lastFetch time.Time
	apiKey    string 
}

func NewConverter(apiKey string) *Converter {
	return &Converter{
		cache:  make(map[string]float64),
		apiKey: apiKey, 
	}
}

func (c *Converter) ConvertToEUR(ctx context.Context, amount int, currency string) (int, error) {
	if currency == "EUR" {
		return amount, nil
	}

	rate, err := c.getRate(ctx, currency)
	if err != nil {
		return 0, err
	}

	return int(math.Round(float64(amount) * rate)), nil
}

func (c *Converter) getRate(ctx context.Context, currency string) (float64, error) {
	c.cacheLock.RLock()
	rate, ok := c.cache[currency]
	cacheValid := time.Since(c.lastFetch) < time.Minute
	c.cacheLock.RUnlock()

	if ok && cacheValid {
		return rate, nil
	}

	rates, err := c.fetchRates(ctx, currency) 
	if err != nil {
		if !ok { 
			return 0, err
		}
		
		return rate, nil
	}

	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()

	for curr, r := range rates {
		c.cache[curr] = r
	}
	c.lastFetch = time.Now()

	if rate, ok := rates[currency]; ok {
		return rate, nil
	}
	return 0, fmt.Errorf("currency not found: %s", currency)
}

func (c *Converter) fetchRates(ctx context.Context, currency string) (map[string]float64, error) {
	url := fmt.Sprintf("https://api.apilayer.com/exchangerates_data/latest?symbols=%s&base=EUR&apikey=%s", currency, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch exchange rates")
	}

	var result struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Convert
	rates := make(map[string]float64)
	for curr, rate := range result.Rates {
		if rate == 0 {
			continue
		}
		rates[curr] = 1 / rate
	}
	return rates, nil
}