package pgsqlrepo

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
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

func (r *expenseLimitsRepo) getDefaultLimit(userId int64, index int) *expenses.ExpenseLimit {
	return &expenses.ExpenseLimit{
		UserId: userId,
		Month:  index,
		Value:  r.service.GetConfig().ExpenseLimitDefault}
}

func (r *expenseLimitsRepo) GetOne(ctx context.Context, userId int64, index int) (*expenses.ExpenseLimit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get expense limit from database")
	defer span.Finish()

	var limit expenses.ExpenseLimit
	err := r.pool.QueryRow(ctx, "select user_id, month, value from expense_limits where user_id = $1 and month = $2", userId, index).
		Scan(&limit.UserId, &limit.Month, &limit.Value)
	if err == pgx.ErrNoRows {
		return r.getDefaultLimit(userId, index), nil
	}
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

func (r *expenseLimitsRepo) GetAll(ctx context.Context, userId int64) ([]expenses.ExpenseLimit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get all expense limits from database")
	defer span.Finish()

	result := make([]expenses.ExpenseLimit, 12)
	rows, err := r.pool.Query(ctx, "select user_id, month, value from expense_limits where user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	limitMap := make(map[int]expenses.ExpenseLimit)
	for rows.Next() {
		var limit expenses.ExpenseLimit
		err := rows.Scan(&limit.UserId, &limit.Month, &limit.Value)
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
			result[i] = *r.getDefaultLimit(userId, i+1)
		}
	}
	return result, nil
}

func (r *expenseLimitsRepo) Save(ctx context.Context, el *expenses.ExpenseLimit) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "save user limit to database")
	defer span.Finish()

	_, err := r.pool.Exec(ctx, `
		insert into expense_limits(user_id, month, value) values($1, $2, $3)
		on conflict (user_id, month) do update set value=excluded.value`,
		el.UserId, el.Month, el.Value)
	return err
}

func (r *expenseLimitsRepo) Delete(ctx context.Context, userId int64, index int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "delete user limit from expense")
	defer span.Finish()

	res, err := r.pool.Exec(ctx, "delete from expense_limits where user_id=$1 and month=$2", userId, index)
	if err == nil {
		if res.RowsAffected() == 0 {
			return localerr.ErrExpenseLimitNotFound
		}
	}
	return err
}
