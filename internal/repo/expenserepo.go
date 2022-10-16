package repo

import (
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
)

type ExpCurrencyDayData map[time.Time]int

type ExpCurrencyData map[string]ExpCurrencyDayData

type ExpData map[string]ExpCurrencyData

type ExpensesRepo interface {
	Add(e *expenses.Expense) error
	ExpensesByUserAndTimeInterval(UserID int64, timeStart time.Time, timeEnd time.Time) (ExpData, error)
}

type ExpenseLimitsRepo interface {
	GetOne(UserID int64, index int) (*expenses.ExpenseLimit, error)
	GetAll(UserID int64) ([]expenses.ExpenseLimit, error)
	Save(e *expenses.ExpenseLimit) error
	Delete(UserID int64, index int) error
}
