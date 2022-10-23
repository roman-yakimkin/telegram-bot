package pgsqlrepo

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type userStateRepo struct {
	pool *pgxpool.Pool
}

func NewUserStateRepo(pool *pgxpool.Pool) repo.UserStateRepo {
	return &userStateRepo{
		pool: pool,
	}
}

func (r *userStateRepo) GetOne(ctx context.Context, UserId int64) (*userstates.UserState, error) {
	var currency string
	var status int
	var rawJson string
	err := r.pool.QueryRow(ctx, "select currency_code, status, input_buffer from user_states where user_id = $1", UserId).
		Scan(&currency, &status, &rawJson)
	if err == pgx.ErrNoRows {
		return nil, localerr.ErrUserStateNotFound
	}
	if err != nil {
		return nil, err
	}
	buffer := make(map[string]interface{})
	err = json.Unmarshal([]byte(rawJson), &buffer)
	if err != nil {
		return nil, err
	}
	return userstates.CreateUserState(UserId, currency, status, buffer), nil
}

func (r *userStateRepo) Save(ctx context.Context, state *userstates.UserState) error {
	state.BeforeSave()
	jsonBuffer, err := state.GetJSONBuffer()
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		insert into user_states (user_id, currency_code, status, input_buffer) values($1, $2, $3, $4) 
		on conflict (user_id) do update set currency_code=excluded.currency_code, status=excluded.status, input_buffer=excluded.input_buffer`,
		state.UserID, state.Currency, state.GetStatus(), jsonBuffer)
	return err
}

func (r *userStateRepo) Delete(ctx context.Context, UserID int64) error {
	res, err := r.pool.Exec(ctx, "delete from user_states where user_id=$1", UserID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return localerr.ErrUserStateNotFound
	}
	return nil
}

func (r *userStateRepo) ClearStatus(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, "update user_states set status = $1", userstates.ExpectedCommand)
	return err
}
