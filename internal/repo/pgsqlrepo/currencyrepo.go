package pgsqlrepo

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type currencyRepo struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func NewCurrencyRepo(ctx context.Context, pool *pgxpool.Pool) repo.CurrencyRepo {
	return &currencyRepo{
		ctx:  ctx,
		pool: pool,
	}
}

func (r *currencyRepo) GetOne(currName string) (*currencies.Currency, error) {
	var curr currencies.Currency
	err := r.pool.QueryRow(r.ctx,
		"select id, display from currency where id=$1", currName).
		Scan(&curr.Name, &curr.Display)
	if err != nil {
		return nil, err
	}
	return &curr, nil
}

func (r *currencyRepo) GetAll() ([]currencies.Currency, error) {
	var result []currencies.Currency
	rows, err := r.pool.Query(r.ctx, "select id, display from currency")
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
