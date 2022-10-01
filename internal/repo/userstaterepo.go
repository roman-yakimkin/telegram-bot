package repo

import "gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"

type UserStateRepo interface {
	GetOne(UserId int64) (*userstates.UserState, error)
	Save(state *userstates.UserState) error
	Delete(UserID int64) error
}
