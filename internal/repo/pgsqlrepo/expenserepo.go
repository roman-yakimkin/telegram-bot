package pgsqlrepo

import (
	"context"
	"log"
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
	pool    *pgxpool.Pool
	service *config.Service
}

func NewExpenseRepo(pool *pgxpool.Pool, service *config.Service) repo.ExpensesRepo {
	return &expensesRepo{
		pool:    pool,
		service: service,
	}
}

func (r *expensesRepo) Add(ctx context.Context, e *expenses.Expense, limitChecker repo.ExpenseLimitChecker) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	catId, err := r.getCategoryId(ctx, tx, e.Category)
	if err != nil {
		if rErr := tx.Rollback(ctx); rErr != nil {
			log.Println("rollback error", rErr)
		}
		return err
	}
	_, err = tx.Exec(ctx, "insert into expenses (user_id, category_id, currency_code, amount, date) values ($1, $2, $3, $4, $5)",
		e.UserID, catId, e.Currency, e.Amount, utils.TimeTruncate(e.Date))
	if err != nil {
		if rErr := tx.Rollback(ctx); rErr != nil {
			log.Println("rollback error", rErr)
		}
		return err
	}
	ok, err := limitChecker.MeetMonthlyLimit(ctx, e.UserID, utils.TimeTruncate(e.Date), e.Amount, limitChecker.CurrencyConvertorTo())
	if !ok || err != nil {
		if !ok {
			log.Println("monthly limit exceeded")
		}
		if rErr := tx.Rollback(ctx); rErr != nil {
			log.Println("rollback error", rErr)
		}
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *expensesRepo) getCategoryId(ctx context.Context, tx pgx.Tx, catName string) (int, error) {
	var id int
	var name string
	err := tx.QueryRow(ctx, "select id, name from categories where upper(name)=$1", strings.ToUpper(catName)).
		Scan(&id, &name)
	if err == pgx.ErrNoRows {
		err := tx.QueryRow(ctx, "insert into categories (name) values ($1) returning id", catName).
			Scan(&id)
		return id, err
	}
	return id, err
}

func (r *expensesRepo) ExpensesByUserAndTimeInterval(ctx context.Context, UserID int64, timeStart time.Time, timeEnd time.Time) (repo.ExpData, error) {
	result := make(repo.ExpData)
	rows, err := r.pool.Query(ctx, `
			select c.name as category_name, e.currency_code, e.amount, e.date 
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
		if _, ok := result[catName]; !ok {
			result[catName] = make(repo.ExpCurrencyData)
		}
		if _, ok := result[catName][curr]; !ok {
			result[catName][curr] = make(repo.ExpCurrencyDayData)
		}
		result[catName][curr][date] += amount
	}
	return result, nil
}
