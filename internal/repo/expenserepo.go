package repo

import (
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
)

type ExpCurrencyData map[string]int
type ExpData map[string]ExpCurrencyData

type ExpensesRepo interface {
	Add(e *expenses.Expense) error
	ExpensesByUserAndTimeInterval(UserID int64, timeStart time.Time, timeEnd time.Time) ExpData
}
