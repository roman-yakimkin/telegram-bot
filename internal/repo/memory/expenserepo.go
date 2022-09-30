package memory

import (
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
)

type ExpensesUserCatPayment struct {
	amount int
	date   time.Time
}

type ExpensesUserCat map[string][]ExpensesUserCatPayment

type ExpenseRepo struct {
	e map[int64]ExpensesUserCat
}

func NewExpenseRepo() *ExpenseRepo {
	return &ExpenseRepo{
		e: make(map[int64]ExpensesUserCat),
	}
}

func (r *ExpenseRepo) Add(e *expenses.Expense) error {
	_, ok := r.e[e.UserID]
	if !ok {
		r.e[e.UserID] = make(ExpensesUserCat)
	}
	payments := r.e[e.UserID][e.Category]
	payments = append(payments, ExpensesUserCatPayment{
		amount: e.Amount,
		date:   e.Date,
	})
	r.e[e.UserID][e.Category] = payments
	return nil
}

func (r *ExpenseRepo) ExpensesByUserAndTimeInterval(UserID int64, timeStart time.Time, timeEnd time.Time) map[string]int {
	result := make(map[string]int)
	userData, ok := r.e[UserID]
	if !ok {
		return result
	}
	for category, payments := range userData {
		sum := 0
		for _, payment := range payments {
			if payment.date.After(timeStart) && payment.date.Before(timeEnd) {
				sum += payment.amount
			}
		}
		if sum > 0 {
			result[category] = sum
		}
	}
	return result
}
