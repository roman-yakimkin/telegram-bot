package pgsqlrepo

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type ExpenseUserLimits = map[int]expenses.ExpenseLimit

type expenseLimitsRepo struct {
	ctx     context.Context
	pool    *pgxpool.Pool
	service *config.Service
}

func NewExpenseLimitsRepo(ctx context.Context, pool *pgxpool.Pool, service *config.Service) repo.ExpenseLimitsRepo {
	return &expenseLimitsRepo{
		ctx:     ctx,
		pool:    pool,
		service: service,
	}
}

func (r *expenseLimitsRepo) GetOne(UserID int64, index int) (*expenses.ExpenseLimit, error) {
	var limit expenses.ExpenseLimit
	err := r.pool.QueryRow(r.ctx, "select user_id, month, value from expense_limits where user_id = $1 and month = $2", UserID, index).
		Scan(&limit.UserID, &limit.Month, &limit.Value)
	if err == nil {
		return &limit, nil
	} else if err == pgx.ErrNoRows {
		return &expenses.ExpenseLimit{
			UserID: UserID,
			Month:  index,
			Value:  r.service.GetConfig().ExpenseLimitDefault}, nil
	} else {
		return nil, err
	}
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
	_, err := r.pool.Exec(r.ctx, `
		insert into expense_limits(user_id, month, value) values($1, $2, $3)
		on conflict (user_id, month) do update set value=excluded.value`,
		el.UserID, el.Month, el.Value)
	return err
}

func (r *expenseLimitsRepo) Delete(UserID int64, index int) error {
	res, err := r.pool.Exec(r.ctx, "delete from expense_limits where user_id=$1 and month=$2", UserID, index)
	if err == nil {
		if res.RowsAffected() == 0 {
			return localerr.ErrExpenseLimitNotFound
		}
	}
	return err
}
