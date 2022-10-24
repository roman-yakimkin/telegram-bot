package importers

import (
	"context"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
	"time"
)

type CurrencyRateImporter interface {
	GetRatesByDate(ctx context.Context, date time.Time) ([]currencies.CurrencyRate, error)
	GetCurrencyCount() int
}
