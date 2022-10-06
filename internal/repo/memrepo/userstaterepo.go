package memrepo

import (
	"sync"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type UserStateRepo struct {
	mx     sync.Mutex
	states map[int64]userstates.UserState
}

func NewUserStateRepo() *UserStateRepo {
	return &UserStateRepo{
		states: make(map[int64]userstates.UserState),
	}
}

func (r *UserStateRepo) GetOne(UserId int64) (*userstates.UserState, error) {
	r.mx.Lock()
	userState, ok := r.states[UserId]
	r.mx.Unlock()
	if !ok {
		return nil, localerr.ErrUserStateNotFound
	}
	return &userState, nil
}

func (r *UserStateRepo) Save(state *userstates.UserState) error {
	state.BeforeSave()
	r.mx.Lock()
	r.states[state.UserID] = *state
	r.mx.Unlock()
	return nil
}

func (r *UserStateRepo) Delete(UserID int64) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	_, ok := r.states[UserID]
	if !ok {
		return localerr.ErrUserStateNotFound
	}
	delete(r.states, UserID)
	return nil
}
