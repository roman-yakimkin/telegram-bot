package repo

import (
	"context"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
)

type ExpCurrencyDayData map[time.Time]int

type ExpCurrencyData map[string]ExpCurrencyDayData

type ExpData map[string]ExpCurrencyData

type CurrencyConvertorFrom interface {
	From(ctx context.Context, amount int, currFrom string, date time.Time) (int, error)
}

type CurrencyConvertorTo interface {
	To(ctx context.Context, amount int, currTo string, date time.Time) (int, error)
}

type CurrencyConvertor interface {
	CurrencyConvertorFrom
	CurrencyConvertorTo
}

type ExpenseLimitChecker interface {
	CurrencyConvertorTo() CurrencyConvertorTo
	MeetMonthlyLimit(ctx context.Context, UserID int64, date time.Time, amountInRub int, curr CurrencyConvertorTo) (bool, error)
}

type ExpensesRepo interface {
	Add(ctx context.Context, e *expenses.Expense, limitChecker ExpenseLimitChecker) error
	ExpensesByUserAndTimeInterval(ctx context.Context, UserID int64, timeStart time.Time, timeEnd time.Time) (ExpData, error)
}

type ExpenseLimitsRepo interface {
	GetOne(ctx context.Context, UserID int64, index int) (*expenses.ExpenseLimit, error)
	GetAll(ctx context.Context, UserID int64) ([]expenses.ExpenseLimit, error)
	Save(ctx context.Context, e *expenses.ExpenseLimit) error
	Delete(ctx context.Context, UserID int64, index int) error
}
