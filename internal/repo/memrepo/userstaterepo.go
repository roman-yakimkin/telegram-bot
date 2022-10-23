package memrepo

import (
	"context"
	"sync"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type userStateRepo struct {
	mx     sync.Mutex
	states map[int64]userstates.UserState
}

func NewUserStateRepo() repo.UserStateRepo {
	return &userStateRepo{
		states: make(map[int64]userstates.UserState),
	}
}

func (r *userStateRepo) GetOne(_ context.Context, UserId int64) (*userstates.UserState, error) {
	r.mx.Lock()
	userState, ok := r.states[UserId]
	r.mx.Unlock()
	if !ok {
		return nil, localerr.ErrUserStateNotFound
	}
	return &userState, nil
}

func (r *userStateRepo) Save(_ context.Context, state *userstates.UserState) error {
	state.BeforeSave()
	r.mx.Lock()
	r.states[state.UserID] = *state
	r.mx.Unlock()
	return nil
}

func (r *userStateRepo) Delete(_ context.Context, UserID int64) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	_, ok := r.states[UserID]
	if !ok {
		return localerr.ErrUserStateNotFound
	}
	delete(r.states, UserID)
	return nil
}

func (r *userStateRepo) ClearStatus(_ context.Context) error {
	return nil
}
