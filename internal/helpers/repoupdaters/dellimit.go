package repoupdaters

import (
	"log"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type delLimitSaver struct {
	userState *userstates.UserState
	limitRepo repo.ExpenseLimitsRepo
}

func NewDelLimitSaver(limitRepo repo.ExpenseLimitsRepo) UserStateRepoUpdater {
	return &delLimitSaver{
		limitRepo: limitRepo,
	}
}

func (s *delLimitSaver) SetUserState(userState *userstates.UserState) {
	s.userState = userState
}

func (s *delLimitSaver) toLimitMonth() (int, error) {
	index, err := s.userState.IfFloatTransformToInt(userstates.DeleteLimitMonthIndex)
	if err != nil {
		log.Println("error upon getting limit month index: ", err)
	}
	return index, err
}

func (s *delLimitSaver) ReadyToUpdate() bool {
	return s.userState.BufferValueExists(userstates.DeleteLimitMonthIndex)
}

func (s *delLimitSaver) UpdateRepo() error {
	index, err := s.toLimitMonth()
	if err != nil {
		return err
	}
	return s.limitRepo.Delete(s.userState.UserID, index)
}

func (s *delLimitSaver) ClearData() {
	s.userState.ClearBufferValue(userstates.DeleteLimitMonthIndex)
}
