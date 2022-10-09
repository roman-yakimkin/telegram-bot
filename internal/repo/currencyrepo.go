package repo

import "gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"

type CurrencyRepo interface {
	GetOne(currName string) (*currencies.Currency, error)
	GetAll() ([]currencies.Currency, error)
}
