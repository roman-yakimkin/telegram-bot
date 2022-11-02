package importers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/utils"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
)

type RawCurrencyData struct {
	Name    string  `json:"CharCode"`
	Nominal int     `json:"Nominal"`
	Value   float64 `json:"Value"`
}

func (rd *RawCurrencyData) toCurrencyRate(date time.Time) currencies.CurrencyRate {
	return currencies.CurrencyRate{
		Name:       rd.Name,
		RateToMain: rd.Value / float64(rd.Nominal),
		Date:       utils.TimeTruncate(date),
	}
}

type cbrRateImporter struct {
	pool    *pgxpool.Pool
	service *config.Service
}

func NewCbrRateImporter(pool *pgxpool.Pool, service *config.Service) CurrencyRateImporter {
	return &cbrRateImporter{
		pool:    pool,
		service: service,
	}
}

func (ri *cbrRateImporter) currencyURL(date time.Time) string {
	cfg := ri.service.GetConfig()
	if date == utils.TimeTruncate(time.Now()) {
		return cfg.CurrencyURLCurrent
	} else {
		y, m, d := date.Date()
		r := fmt.Sprintf(cfg.CurrencyURLPast, y, m, d)
		return r
	}
}

func (ri *cbrRateImporter) isRateUnset(resp *http.Response) bool {
	cfg := ri.service.GetConfig()
	var rawData struct {
		Explanation string `json:"explanation"`
	}
	err := json.NewDecoder(resp.Body).Decode(&rawData)
	return err == nil && strings.Contains(rawData.Explanation, cfg.CurrencyRateUnset)
}

func (ri *cbrRateImporter) GetRatesByDate(ctx context.Context, date time.Time) ([]currencies.CurrencyRate, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "import currency rate from cbi service")
	defer span.Finish()

	cfg := ri.service.GetConfig()
	resp, err := http.Get(ri.currencyURL(date))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound && ri.isRateUnset(resp) {
		return nil, localerr.ErrCurrencyRateUnset
	}
	if resp.StatusCode != http.StatusOK {
		return nil, localerr.ErrCannotGetCurrencyRates
	}
	var rawData struct {
		Currencies map[string]RawCurrencyData `json:"Valute"`
	}
	rawData.Currencies = make(map[string]RawCurrencyData)
	err = json.NewDecoder(resp.Body).Decode(&rawData)
	if err != nil {
		return nil, err
	}
	currRates := make([]currencies.CurrencyRate, 0, len(cfg.Currencies)+1)
	rawData.Currencies["RUB"] = RawCurrencyData{
		Name:    "RUB",
		Nominal: 1,
		Value:   1,
	}
	for _, cfgCurr := range cfg.Currencies {
		rawCurr := rawData.Currencies[cfgCurr.Name]
		currRates = append(currRates, rawCurr.toCurrencyRate(date))
	}
	return currRates, nil
}

func (ri *cbrRateImporter) GetCurrencyCount() int {
	return len(ri.service.GetConfig().Currencies)
}
