package memrepo

import (
	"sync"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type ExpensesUserCatPayment struct {
	amount   int
	currency string
	date     time.Time
}

type ExpensesUserCat map[string][]ExpensesUserCatPayment

type ExpenseRepo struct {
	mx sync.Mutex
	e  map[int64]ExpensesUserCat
}

func NewExpenseRepo() *ExpenseRepo {
	return &ExpenseRepo{
		e: make(map[int64]ExpensesUserCat),
	}
}

func (r *ExpenseRepo) Add(e *expenses.Expense) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	_, ok := r.e[e.UserID]
	if !ok {
		r.e[e.UserID] = make(ExpensesUserCat)
	}
	payments := r.e[e.UserID][e.Category]
	payments = append(payments, ExpensesUserCatPayment{
		amount:   e.Amount,
		currency: e.Currency,
		date:     e.Date,
	})
	r.e[e.UserID][e.Category] = payments
	return nil
}

func (r *ExpenseRepo) ExpensesByUserAndTimeInterval(UserID int64, timeStart time.Time, timeEnd time.Time) repo.ExpData {
	r.mx.Lock()
	defer r.mx.Unlock()
	result := make(repo.ExpData)
	userData, ok := r.e[UserID]
	if !ok {
		return result
	}
	for category, payments := range userData {
		sums := make(repo.ExpCurrencyData)
		for _, payment := range payments {
			if payment.date.After(timeStart) && payment.date.Before(timeEnd) {
				sums[payment.currency] += payment.amount
			}
		}
		if len(sums) > 0 {
			result[category] = sums
		}
	}
	return result
}
