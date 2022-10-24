package memrepo

import (
	"context"
	"sync"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/utils"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type ExpensesUserCatPayment struct {
	amount   int
	currency string
	date     time.Time
}

type ExpensesUserCat map[string][]ExpensesUserCatPayment

type expensesRepo struct {
	mx sync.Mutex
	e  map[int64]ExpensesUserCat
}

func NewExpenseRepo() repo.ExpensesRepo {
	return &expensesRepo{
		e: make(map[int64]ExpensesUserCat),
	}
}

func (r *expensesRepo) Add(_ context.Context, e *expenses.Expense, _ repo.ExpenseLimitChecker) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	_, ok := r.e[e.UserId]
	if !ok {
		r.e[e.UserId] = make(ExpensesUserCat)
	}
	payments := r.e[e.UserId][e.Category]
	payments = append(payments, ExpensesUserCatPayment{
		amount:   e.Amount,
		currency: e.Currency,
		date:     utils.TimeTruncate(e.Date),
	})
	r.e[e.UserId][e.Category] = payments
	return nil
}

func (r *expensesRepo) ExpensesByUserAndTimeInterval(_ context.Context, userId int64, timeStart time.Time, timeEnd time.Time) (repo.ExpData, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	result := make(repo.ExpData)
	userData, ok := r.e[userId]
	if !ok {
		return result, nil
	}
	for category, payments := range userData {
		sums := make(repo.ExpCurrencyData)
		for _, payment := range payments {
			if payment.date.After(timeStart) && payment.date.Before(timeEnd) {
				mapData := sums[payment.currency]
				if mapData == nil {
					mapData = make(repo.ExpCurrencyDayData)
				}
				mapData[utils.TimeTruncate(payment.date)] += payment.amount
				sums[payment.currency] = mapData
			}
		}
		if len(sums) > 0 {
			result[category] = sums
		}
	}
	return result, nil
}
