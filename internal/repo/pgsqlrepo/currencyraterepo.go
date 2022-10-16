package pgsqlrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/utils"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type сurrencyRateRepo struct {
	mx      sync.Mutex
	ctx     context.Context
	pool    *pgxpool.Pool
	service *config.Service
}

func NewCurrencyRateRepo(ctx context.Context, pool *pgxpool.Pool, service *config.Service) repo.CurrencyRateRepo {
	return &сurrencyRateRepo{
		ctx:     ctx,
		pool:    pool,
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
	Nominal int     `json:"Nominal"`
	Value   float64 `json:"Value"`
}

func (r *сurrencyRateRepo) isRateUnset(resp *http.Response) bool {
	cfg := r.service.GetConfig()
	var rawData struct {
		Explanation string `json:"explanation"`
	}
	err := json.NewDecoder(resp.Body).Decode(&rawData)
	return err == nil && strings.Contains(rawData.Explanation, cfg.CurrencyRateUnset)
}

func (r *сurrencyRateRepo) LoadByDate(date time.Time) error {
	cfg := r.service.GetConfig()
	resp, err := http.Get(r.currencyURL(date))
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNotFound && r.isRateUnset(resp) {
		return localerr.ErrCurrencyRateUnset
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
		_, err := r.pool.Exec(r.ctx,
			`insert into currency_rates(currency_id, date, rate) values ($1, $2, $3)
				on conflict (currency_id, date) do update set rate=excluded.rate`,
			curr.Name, utils.TimeTruncate(date), thisCurr.Value/float64(thisCurr.Nominal))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *сurrencyRateRepo) LoadByDateIfEmpty(date time.Time) error {
	has, err := r.HasRatesByDate(date)
	if err != nil {
		return err
	}
	if !has {
		err := r.LoadByDate(date)
		if err == localerr.ErrCurrencyRateUnset {
			return r.LoadByDateIfEmpty(date.AddDate(0, 0, -1))
		}
		return err
	}
	return nil
}

func (r *сurrencyRateRepo) loadByDateRecursive(date time.Time) error {
	err := r.LoadByDate(date)
	if err == localerr.ErrCurrencyRateUnset {
		return r.LoadByDate(date.AddDate(0, 0, -1))
	}
	return err
}

func (r *сurrencyRateRepo) GetOneByDate(currName string, date time.Time) (*currencies.CurrencyRate, error) {
	var currRate currencies.CurrencyRate
	for currRate.Name == "" {
		err := r.pool.QueryRow(r.ctx, "select currency_id, date, rate from currency_rates where currency_id = $1 and date = $2", currName, date).
			Scan(&currRate.Name, &currRate.Date, &currRate.RateToMain)
		if err == nil {
			continue
		}
		if err == pgx.ErrNoRows {
			err := r.loadByDateRecursive(date)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &currRate, nil
}

func (r *сurrencyRateRepo) HasRatesByDate(date time.Time) (bool, error) {
	var cntRows int
	err := r.pool.QueryRow(r.ctx, "select count(currency_id) from currency_rates where date = $1", date).Scan(&cntRows)
	if err != nil {
		return false, err
	}
	return cntRows == len(r.service.GetConfig().Currencies), nil
}

func (r *сurrencyRateRepo) GetAllByDate(date time.Time) ([]currencies.CurrencyRate, error) {
	err := r.LoadByDateIfEmpty(date)
	if err != nil {
		return nil, err
	}
	rows, err := r.pool.Query(r.ctx, "select currency_id, date, rate from currency_rates where date = $1", date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []currencies.CurrencyRate
	for rows.Next() {
		var cr currencies.CurrencyRate
		err := rows.Scan(&cr.Name, &cr.Date, &cr.RateToMain)
		if err != nil {
			return nil, err
		}
		result = append(result, cr)
	}
	return result, nil
}

func (r *сurrencyRateRepo) GetAll() ([]currencies.CurrencyRate, error) {

	rows, err := r.pool.Query(r.ctx, "select currency_id, date, rate from currency_rates order by date")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []currencies.CurrencyRate
	for rows.Next() {
		var cr currencies.CurrencyRate
		err := rows.Scan(&cr.Name, &cr.Date, &cr.RateToMain)
		if err != nil {
			return nil, err
		}
		result = append(result, cr)
	}
	return result, nil
}

//func (r *сurrencyRateRepo) dateMapToSlice(currencyMap map[string]currencies.CurrencyRate) []currencies.CurrencyRate {
//	result := make([]currencies.CurrencyRate, 0, len(currencyMap))
//	for _, currency := range currencyMap {
//		result = append(result, currency)
//	}
//	return result
//}
