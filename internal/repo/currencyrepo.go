package repo

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
)

type CurrencyRepo interface {
	GetOne(ctx context.Context, currName string) (*currencies.Currency, error)
	GetAll(ctx context.Context) ([]currencies.Currency, error)
}
