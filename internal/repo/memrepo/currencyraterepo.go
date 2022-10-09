package memrepo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/utils"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type сurrencyRateRepo struct {
	mx      sync.Mutex
	c       map[time.Time]map[string]currencies.CurrencyRate
	service *config.Service
}

func NewCurrencyRateRepo(service *config.Service) repo.CurrencyRateRepo {
	return &сurrencyRateRepo{
		c:       make(map[time.Time]map[string]currencies.CurrencyRate),
		service: service,
	}
}

func (r *сurrencyRateRepo) currencyURL(date time.Time) string {
	cfg := r.service.GetConfig()
	if date == utils.TimeTruncate(time.Now()) {
		return cfg.CurrencyURLCurrent
	} else {
		y, m, d := date.Date()
		r := fmt.Sprintf(cfg.CurrencyURLPast, y, m, d)
		return r
	}
}

type RawCurrencyData struct {
	Name    string  `json:"CharCode"`
	Nominal int     `json:"nominal"`
	Value   float64 `json:"Value"`
}

func (r *сurrencyRateRepo) LoadByDate(date time.Time) error {
	cfg := r.service.GetConfig()
	resp, err := http.Get(r.currencyURL(date))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return localerr.ErrCannotGetCurrencyRates
	}
	var rawData struct {
		Currencies map[string]RawCurrencyData `json:"Valute"`
	}
	rawData.Currencies = make(map[string]RawCurrencyData)
	err = json.NewDecoder(resp.Body).Decode(&rawData)
	if err != nil {
		return err
	}
	r.mx.Lock()
	currMap := make(map[string]currencies.CurrencyRate, len(cfg.Currencies))
	rawData.Currencies["RUB"] = RawCurrencyData{
		Name:    "RUB",
		Nominal: 1,
		Value:   1,
	}

	for _, curr := range cfg.Currencies {
		thisCurr, ok := rawData.Currencies[curr.Name]
		if !ok {
			return localerr.ErrCurrencyNotFound
		}
		currMap[curr.Name] = currencies.CurrencyRate{
			Name:       curr.Name,
			RateToMain: thisCurr.Value / float64(thisCurr.Nominal),
			Date:       date,
		}
	}
	r.c[date] = currMap
	r.mx.Unlock()
	return nil
}

func (r *сurrencyRateRepo) GetOneByDate(currName string, date time.Time) (*currencies.CurrencyRate, error) {
	r.mx.Lock()
	currenciesByDate, ok := r.c[date]
	r.mx.Unlock()
	if !ok {
		err := r.LoadByDate(date)
		if err != nil {
			return nil, err
		}
		currenciesByDate = r.c[date]
	}
	currency, ok := currenciesByDate[currName]
	if !ok {
		return nil, localerr.ErrCurrencyNotFound
	}
	return &currency, nil
}

func (r *сurrencyRateRepo) GetAllByDate(date time.Time) ([]currencies.CurrencyRate, error) {
	r.mx.Lock()
	currenciesByDate, ok := r.c[date]
	r.mx.Unlock()
	if !ok {
		err := r.LoadByDate(date)
		if err != nil {
			return nil, err
		}
		r.mx.Lock()
		currenciesByDate = r.c[date]
		r.mx.Unlock()
	}
	return r.dateMapToSlice(currenciesByDate), nil
}

func (r *сurrencyRateRepo) GetAll() ([]currencies.CurrencyRate, error) {
	var result []currencies.CurrencyRate

	for date := range r.c {
		currDateSlice, err := r.GetAllByDate(date)
		if err != nil {
			return nil, err
		}
		result = append(result, currDateSlice...)
	}
	return result, nil
}

func (r *сurrencyRateRepo) dateMapToSlice(currencyMap map[string]currencies.CurrencyRate) []currencies.CurrencyRate {
	result := make([]currencies.CurrencyRate, 0, len(currencyMap))
	for _, currency := range currencyMap {
		result = append(result, currency)
	}
	return result
}
