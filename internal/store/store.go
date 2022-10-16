package store

import (
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type UserCurrencyGetter interface {
	UserCurrencyRate(UserID int64, date time.Time) (*currencies.CurrencyRate, error)
}

type Store interface {
	UserCurrencyGetter
	CurrencyRate() repo.CurrencyRateRepo
	Currency() repo.CurrencyRepo
	UserState() repo.UserStateRepo
	Expense() repo.ExpensesRepo
	Limit() repo.ExpenseLimitsRepo
	MeetMonthlyLimit(UserID int64, date time.Time, amountInRub int, curr convertors.CurrencyConvertorTo) (bool, error)
}
