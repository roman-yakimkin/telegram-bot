package repoupdaters

import (
	"log"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type saveLimitSaver struct {
	userState *userstates.UserState
	limitRepo repo.ExpenseLimitsRepo
}

func NewSaveLimitSaver(limitRepo repo.ExpenseLimitsRepo) UserStateRepoUpdater {
	return &saveLimitSaver{
		limitRepo: limitRepo,
	}
}

func (s *saveLimitSaver) SetUserState(userState *userstates.UserState) {
	s.userState = userState
}

func (s *saveLimitSaver) toLimit() (*expenses.ExpenseLimit, error) {
	month, err := s.userState.IfFloatTransformToInt(userstates.SetLimitMonthIndex)
	if err != nil {
		log.Println("error upon getting month index", err)
		return nil, err
	}
	value, err := s.userState.IfFloatTransformToInt(userstates.SetLimitMonthValue)
	if err != nil {
		log.Println("error upon getting month value", err)
		return nil, err
	}
	return &expenses.ExpenseLimit{
		UserID: s.userState.UserID,
		Month:  month,
		Value:  value,
	}, nil
}

func (s *saveLimitSaver) ReadyToUpdate() bool {
	ok1 := s.userState.BufferValueExists(userstates.SetLimitMonthIndex)
	ok2 := s.userState.BufferValueExists(userstates.SetLimitMonthValue)
	return ok1 && ok2
}

func (s *saveLimitSaver) UpdateRepo() error {
	limit, err := s.toLimit()
	if err != nil {
		return err
	}
	return s.limitRepo.Save(limit)
}

func (s *saveLimitSaver) ClearData() {
	s.userState.ClearBufferValue(userstates.SetLimitMonthIndex)
	s.userState.ClearBufferValue(userstates.SetLimitMonthValue)
}
