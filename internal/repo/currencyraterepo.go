package repo

import (
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
)

type CurrencyRateRepo interface {
	LoadByDate(date time.Time) error
	GetOneByDate(currName string, date time.Time) (*currencies.CurrencyRate, error)
	GetAllByDate(date time.Time) ([]currencies.CurrencyRate, error)
	GetAll() ([]currencies.CurrencyRate, error)
}
