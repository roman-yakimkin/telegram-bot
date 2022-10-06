package memrepo

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
)

type CurrencyRepo struct {
	mx      sync.RWMutex
	c       map[string]currencies.Currency
	service *config.Service
}

func NewCurrencyRepo(service *config.Service) *CurrencyRepo {
	return &CurrencyRepo{
		c:       make(map[string]currencies.Currency),
		service: service,
	}
}

func (r *CurrencyRepo) LoadAll() error {
	cfg := r.service.GetConfig()
	resp, err := http.Get(cfg.CurrencyURL)
	if err != nil {
		return err
	}
	var rawData struct {
		Base         string             `json:"base"`
		Rates        map[string]float64 `json:"rates"`
		Source       string             `json:"source"`
		LocalISODate string             `json:"localISODate"`
		PutISODate   string             `json:"putISODate"`
	}
	err = json.NewDecoder(resp.Body).Decode(&rawData)
	if err != nil {
		return err
	}
	r.mx.Lock()
	r.c = make(map[string]currencies.Currency, len(cfg.Currencies))
	for _, curr := range cfg.Currencies {
		r.c[curr.Name] = currencies.Currency{
			Name:       curr.Name,
			Display:    curr.Display,
			RateToMain: rawData.Rates[cfg.CurrencyMain] / rawData.Rates[curr.Name],
			Received:   time.Now(),
		}
	}
	r.mx.Unlock()
	return nil
}

func (r *CurrencyRepo) GetOne(currName string) (*currencies.Currency, error) {
	r.mx.RLock()
	currency, ok := r.c[currName]
	r.mx.RUnlock()
	if !ok {
		return nil, localerr.ErrCurrencyNotFound
	}
	return &currency, nil
}

func (r *CurrencyRepo) GetAll() ([]currencies.Currency, error) {
	result := make([]currencies.Currency, 0, len(r.c))
	r.mx.RLock()
	for _, currency := range r.c {
		result = append(result, currency)
	}
	r.mx.RUnlock()
	return result, nil
}
