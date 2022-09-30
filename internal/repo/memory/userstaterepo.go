package memory

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type UserStateRepo struct {
	states map[int64]userstates.UserState
}

func NewUserStateRepo() *UserStateRepo {
	return &UserStateRepo{
		states: make(map[int64]userstates.UserState),
	}
}

func (r *UserStateRepo) GetOne(UserId int64) (*userstates.UserState, error) {
	userState, ok := r.states[UserId]
	if !ok {
		return nil, localerr.ErrUserStateNotFound
	}
	return &userState, nil
}

func (r *UserStateRepo) Save(state *userstates.UserState) error {
	r.states[state.UserID] = *state
	return nil
}

func (r *UserStateRepo) Delete(UserID int64) error {
	_, ok := r.states[UserID]
	if !ok {
		return localerr.ErrUserStateNotFound
	}
	delete(r.states, UserID)
	return nil
}
