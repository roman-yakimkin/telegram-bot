package store

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type UserCurrencyGetter interface {
	UserCurrency(UserID int64) (*currencies.Currency, error)
}

type Store interface {
	UserCurrencyGetter
	Currency() repo.CurrencyRepo
	UserState() repo.UserStateRepo
	Expense() repo.ExpensesRepo
}
