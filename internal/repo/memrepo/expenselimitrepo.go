package memrepo

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type ExpenseUserLimits = map[int]expenses.ExpenseLimit

type expenseLimitsRepo struct {
	limits  map[int64]ExpenseUserLimits
	service *config.Service
}

func NewExpenseLimitsRepo(service *config.Service) repo.ExpenseLimitsRepo {
	return &expenseLimitsRepo{
		limits:  make(map[int64]ExpenseUserLimits),
		service: service,
	}
}

func (r *expenseLimitsRepo) GetOne(UserID int64, index int) (*expenses.ExpenseLimit, error) {
	limits, ok := r.limits[UserID]
	if !ok {
		return &expenses.ExpenseLimit{
			Month: index,
			Value: r.service.GetConfig().ExpenseLimitDefault}, nil
	}
	limit, ok := limits[index]
	if !ok {
		return &expenses.ExpenseLimit{
			Month: index,
			Value: r.service.GetConfig().ExpenseLimitDefault}, nil
	}
	return &limit, nil
}

func (r *expenseLimitsRepo) GetAll(UserID int64) ([]expenses.ExpenseLimit, error) {
	result := make([]expenses.ExpenseLimit, 12)
	for i := 0; i < len(result); i++ {
		limit, err := r.GetOne(UserID, i+1)
		if err != nil {
			return nil, err
		}
		result[i] = *limit
	}
	return result, nil
}

func (r *expenseLimitsRepo) Save(el *expenses.ExpenseLimit) error {
	_, ok := r.limits[el.UserID]
	if !ok {
		r.limits[el.UserID] = make(ExpenseUserLimits, 12)
	}
	r.limits[el.UserID][el.Month] = *el
	return nil
}

func (r *expenseLimitsRepo) Delete(UserID int64, index int) error {
	limits, ok := r.limits[UserID]
	if !ok {
		return localerr.ErrExpenseLimitNotFound
	}
	_, ok = limits[index]
	if !ok {
		return localerr.ErrExpenseLimitNotFound
	}
	delete(limits, index)
	r.limits[UserID] = limits
	return nil
}
