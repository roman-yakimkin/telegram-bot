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
	pool    *pgxpool.Pool
	service *config.Service
}

func NewExpenseLimitsRepo(pool *pgxpool.Pool, service *config.Service) repo.ExpenseLimitsRepo {
	return &expenseLimitsRepo{
		pool:    pool,
		service: service,
	}
}

func (r *expenseLimitsRepo) getDefaultLimit(UserID int64, index int) *expenses.ExpenseLimit {
	return &expenses.ExpenseLimit{
		UserID: UserID,
		Month:  index,
		Value:  r.service.GetConfig().ExpenseLimitDefault}
}

func (r *expenseLimitsRepo) GetOne(ctx context.Context, UserID int64, index int) (*expenses.ExpenseLimit, error) {
	var limit expenses.ExpenseLimit
	err := r.pool.QueryRow(ctx, "select user_id, month, value from expense_limits where user_id = $1 and month = $2", UserID, index).
		Scan(&limit.UserID, &limit.Month, &limit.Value)
	if err == pgx.ErrNoRows {
		return r.getDefaultLimit(UserID, index), nil
	}
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

func (r *expenseLimitsRepo) GetAll(ctx context.Context, UserID int64) ([]expenses.ExpenseLimit, error) {
	result := make([]expenses.ExpenseLimit, 12)
	rows, err := r.pool.Query(ctx, "select user_id, month, value from expense_limits where user_id = $1", UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	limitMap := make(map[int]expenses.ExpenseLimit)
	for rows.Next() {
		var limit expenses.ExpenseLimit
		err := rows.Scan(&limit.UserID, &limit.Month, &limit.Value)
		if err != nil {
			return nil, err
		}
		limitMap[limit.Month] = limit
	}
	for i := 0; i < len(result); i++ {
		limit, ok := limitMap[i+1]
		if ok {
			result[i] = limit
		} else {
			result[i] = *r.getDefaultLimit(UserID, i+1)
		}
	}
	return result, nil
}

func (r *expenseLimitsRepo) Save(ctx context.Context, el *expenses.ExpenseLimit) error {
	_, err := r.pool.Exec(ctx, `
		insert into expense_limits(user_id, month, value) values($1, $2, $3)
		on conflict (user_id, month) do update set value=excluded.value`,
		el.UserID, el.Month, el.Value)
	return err
}

func (r *expenseLimitsRepo) Delete(ctx context.Context, UserID int64, index int) error {
	res, err := r.pool.Exec(ctx, "delete from expense_limits where user_id=$1 and month=$2", UserID, index)
	if err == nil {
		if res.RowsAffected() == 0 {
			return localerr.ErrExpenseLimitNotFound
		}
	}
	return err
}
