package interfaces

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
	"time"
)

type ExpensesRepo interface {
	Add(e *expenses.Expense) error
	ExpensesByUserAndTimeInterval(UserID int64, timeStart time.Time, timeEnd time.Time) map[string]int
}
