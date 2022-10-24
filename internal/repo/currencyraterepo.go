package repo

import (
	"context"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
)

type CurrencyRateRepo interface {
	LoadByDateIfEmpty(ctx context.Context, date time.Time) error
	GetOneByDate(ctx context.Context, currName string, date time.Time) (*currencies.CurrencyRate, error)
	HasRatesByDate(ctx context.Context, date time.Time) (bool, error)
	GetAllByDate(ctx context.Context, date time.Time) ([]currencies.CurrencyRate, error)
	GetAll(ctx context.Context) ([]currencies.CurrencyRate, error)
}
