package pgsqlrepo

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type currencyRepo struct {
	pool *pgxpool.Pool
}

func NewCurrencyRepo(pool *pgxpool.Pool) repo.CurrencyRepo {
	return &currencyRepo{
		pool: pool,
	}
}

func (r *currencyRepo) GetOne(ctx context.Context, currName string) (*currencies.Currency, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get currency from database")
	defer span.Finish()

	var curr currencies.Currency
	err := r.pool.QueryRow(ctx,
		"select code, display from currency where code=$1", currName).
		Scan(&curr.Name, &curr.Display)
	if err != nil {
		return nil, err
	}
	return &curr, nil
}

func (r *currencyRepo) GetAll(ctx context.Context) ([]currencies.Currency, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get all currencies from database")
	defer span.Finish()

	var result []currencies.Currency
	rows, err := r.pool.Query(ctx, "select code, display from currency")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var curr currencies.Currency
		err := rows.Scan(&curr.Name, &curr.Display)
		if err != nil {
			return nil, err
		}
		result = append(result, curr)
	}
	return result, nil
}
