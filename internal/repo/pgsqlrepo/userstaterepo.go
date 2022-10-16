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
	ctx  context.Context
	pool *pgxpool.Pool
}

func NewUserStateRepo(ctx context.Context, pool *pgxpool.Pool) repo.UserStateRepo {
	return &userStateRepo{
		ctx:  ctx,
		pool: pool,
	}
}

func (r *userStateRepo) GetOne(UserId int64) (*userstates.UserState, error) {
	var currency string
	var status int
	var rawJson string
	err := r.pool.QueryRow(r.ctx, "select currency_id, status, input_buffer from user_states where user_id = $1", UserId).
		Scan(&currency, &status, &rawJson)
	if err == nil {
		buffer := make(map[string]interface{})
		err := json.Unmarshal([]byte(rawJson), &buffer)
		if err != nil {
			return nil, err
		}
		return userstates.CreateUserState(UserId, currency, status, buffer), nil
	}
	if err == pgx.ErrNoRows {
		return nil, localerr.ErrUserStateNotFound
	}
	return nil, err
}

func (r *userStateRepo) Save(state *userstates.UserState) error {
	state.BeforeSave()
	jsonBuffer, err := state.GetJSONBuffer()
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(r.ctx, `
		insert into user_states (user_id, currency_id, status, input_buffer) values($1, $2, $3, $4) 
		on conflict (user_id) do update set currency_id=excluded.currency_id, status=excluded.status, input_buffer=excluded.input_buffer`,
		state.UserID, state.Currency, state.GetStatus(), jsonBuffer)
	return err
}

func (r *userStateRepo) Delete(UserID int64) error {
	res, err := r.pool.Exec(r.ctx, "delete from user_states where user_id=$1", UserID)
	if err == nil {
		if res.RowsAffected() == 0 {
			return localerr.ErrUserStateNotFound
		}
	}
	return err
}
