package pgsqlrepo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/importers"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type сurrencyRateRepo struct {
	pool       *pgxpool.Pool
	currImport importers.CurrencyRateImporter
}

func NewCurrencyRateRepo(pool *pgxpool.Pool, currImport importers.CurrencyRateImporter) repo.CurrencyRateRepo {
	return &сurrencyRateRepo{
		pool:       pool,
		currImport: currImport,
	}
}

func (r *сurrencyRateRepo) LoadByDate(ctx context.Context, date time.Time) error {
	importedCurr, err := r.currImport.GetRatesByDate(ctx, date)
	if err != nil {
		return err
	}
	for _, curr := range importedCurr {
		_, err := r.pool.Exec(ctx,
			`insert into currency_rates(currency_code, date, rate) values ($1, $2, $3)
				on conflict (currency_code, date) do update set rate=excluded.rate`,
			curr.Name, curr.Date, curr.RateToMain)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *сurrencyRateRepo) LoadByDateIfEmpty(ctx context.Context, date time.Time) error {
	has, err := r.HasRatesByDate(ctx, date)
	if err != nil {
		return err
	}
	if !has {
		err := r.LoadByDate(ctx, date)
		if err == localerr.ErrCurrencyRateUnset {
			return r.LoadByDateIfEmpty(ctx, date.AddDate(0, 0, -1))
		}
		return err
	}
	return nil
}

func (r *сurrencyRateRepo) loadByDateRecursive(ctx context.Context, date time.Time) error {
	err := r.LoadByDate(ctx, date)
	if err == localerr.ErrCurrencyRateUnset {
		return r.LoadByDate(ctx, date.AddDate(0, 0, -1))
	}
	return err
}

func (r *сurrencyRateRepo) GetOneByDate(ctx context.Context, currName string, date time.Time) (*currencies.CurrencyRate, error) {
	var currRate currencies.CurrencyRate
	for currRate.Name == "" {
		err := r.pool.QueryRow(ctx, "select currency_code, date, rate from currency_rates where currency_code = $1 and date = $2", currName, date).
			Scan(&currRate.Name, &currRate.Date, &currRate.RateToMain)
		if err == pgx.ErrNoRows {
			err := r.loadByDateRecursive(ctx, date)
			if err != nil {
				return nil, err
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return &currRate, nil
}

func (r *сurrencyRateRepo) HasRatesByDate(ctx context.Context, date time.Time) (bool, error) {
	var cntRows int
	err := r.pool.QueryRow(ctx, "select count(currency_code) from currency_rates where date = $1", date).Scan(&cntRows)
	if err != nil {
		return false, err
	}
	return cntRows == r.currImport.GetCurrencyCount(), nil
}

func (r *сurrencyRateRepo) GetAllByDate(ctx context.Context, date time.Time) ([]currencies.CurrencyRate, error) {
	err := r.LoadByDateIfEmpty(ctx, date)
	if err != nil {
		return nil, err
	}
	rows, err := r.pool.Query(ctx, "select currency_code, date, rate from currency_rates where date = $1", date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []currencies.CurrencyRate
	for rows.Next() {
		var cr currencies.CurrencyRate
		err := rows.Scan(&cr.Name, &cr.Date, &cr.RateToMain)
		if err != nil {
			return nil, err
		}
		result = append(result, cr)
	}
	return result, nil
}

func (r *сurrencyRateRepo) GetAll(ctx context.Context) ([]currencies.CurrencyRate, error) {
	rows, err := r.pool.Query(ctx, "select currency_code, date, rate from currency_rates order by date")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []currencies.CurrencyRate
	for rows.Next() {
		var cr currencies.CurrencyRate
		err := rows.Scan(&cr.Name, &cr.Date, &cr.RateToMain)
		if err != nil {
			return nil, err
		}
		result = append(result, cr)
	}
	return result, nil
}
