package pgsqlrepo

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/utils"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type expensesRepo struct {
	ctx     context.Context
	pool    *pgxpool.Pool
	service *config.Service
}

func NewExpenseRepo(ctx context.Context, pool *pgxpool.Pool, service *config.Service) repo.ExpensesRepo {
	return &expensesRepo{
		ctx:     ctx,
		pool:    pool,
		service: service,
	}
}

func (r *expensesRepo) Add(e *expenses.Expense) error {
	tx, err := r.pool.Begin(r.ctx)
	if err != nil {
		return err
	}
	catId, err := r.getCategoryId(tx, e.Category)
	if err != nil {
		err := tx.Rollback(r.ctx)
		if err != nil {
			return err
		}
		return err
	}
	_, err = tx.Exec(r.ctx, "insert into expenses (user_id, category_id, currency_id, amount, date) values ($1, $2, $3, $4, $5)",
		e.UserID, catId, e.Currency, e.Amount, utils.TimeTruncate(e.Date))
	if err != nil {
		err := tx.Rollback(r.ctx)
		if err != nil {
			return err
		}
		return err
	}
	err = tx.Commit(r.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *expensesRepo) getCategoryId(tx pgx.Tx, catName string) (int, error) {
	var id int
	var name string
	err := tx.QueryRow(r.ctx, "select id, name from categories where upper(name)=$1", strings.ToUpper(catName)).
		Scan(&id, &name)
	if err == nil {
		return id, err
	}
	if err == pgx.ErrNoRows {
		err := tx.QueryRow(r.ctx, "insert into categories (name) values ($1) returning id", catName).
			Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	}
	return 0, err
}

func (r *expensesRepo) ExpensesByUserAndTimeInterval(UserID int64, timeStart time.Time, timeEnd time.Time) (repo.ExpData, error) {
	result := make(repo.ExpData)
	rows, err := r.pool.Query(r.ctx, `
			select c.name as category_name, e.currency_id, e.amount, e.date 
			from expenses e inner join categories c on c.id = e.category_id
			where e.user_id = $1 and e.date between $2 and $3`, UserID, utils.TimeTruncate(timeStart), utils.TimeTruncate(timeEnd))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var catName, curr string
		var amount int
		var date time.Time
		err := rows.Scan(&catName, &curr, &amount, &date)
		if err != nil {
			return nil, err
		}
		if result[catName] == nil {
			result[catName] = make(repo.ExpCurrencyData)
		}
		if result[catName][curr] == nil {
			result[catName][curr] = make(repo.ExpCurrencyDayData)
		}
		result[catName][curr][date] += amount
	}
	return result, nil
}
