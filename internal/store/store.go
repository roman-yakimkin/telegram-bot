package store

import (
	"context"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type UserCurrencyGetter interface {
	UserCurrencyRate(ctx context.Context, UserID int64, date time.Time) (*currencies.CurrencyRate, error)
}

type Store interface {
	UserCurrencyGetter
	CurrencyRate() repo.CurrencyRateRepo
	Currency() repo.CurrencyRepo
	UserState() repo.UserStateRepo
	Expense() repo.ExpensesRepo
	CurrencyConvertorFrom() repo.CurrencyConvertorFrom
	CurrencyConvertorTo() repo.CurrencyConvertorTo
	Limit() repo.ExpenseLimitsRepo

	MeetMonthlyLimit(ctx context.Context, UserID int64, date time.Time, amountInRub int, curr repo.CurrencyConvertorTo) (bool, error)
}
