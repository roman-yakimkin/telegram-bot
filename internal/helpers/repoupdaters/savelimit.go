package repoupdaters

import (
	"context"
	"log"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type saveLimitSaver struct {
	limitRepo repo.ExpenseLimitsRepo
}

func NewSaveLimitSaver(limitRepo repo.ExpenseLimitsRepo) UserStateRepoUpdater {
	return &saveLimitSaver{
		limitRepo: limitRepo,
	}
}

func (s *saveLimitSaver) toLimit(state *userstates.UserState) (*expenses.ExpenseLimit, error) {
	month, err := state.IfFloatTransformToInt(userstates.SetLimitMonthIndex)
	if err != nil {
		log.Println("error upon getting month index", err)
		return nil, err
	}
	value, err := state.IfFloatTransformToInt(userstates.SetLimitMonthValue)
	if err != nil {
		log.Println("error upon getting month value", err)
		return nil, err
	}
	return &expenses.ExpenseLimit{
		UserId: state.UserId,
		Month:  month,
		Value:  value,
	}, nil
}

func (s *saveLimitSaver) ReadyToUpdate(state *userstates.UserState) bool {
	ok1 := state.BufferValueExists(userstates.SetLimitMonthIndex)
	ok2 := state.BufferValueExists(userstates.SetLimitMonthValue)
	return ok1 && ok2
}

func (s *saveLimitSaver) UpdateRepo(ctx context.Context, state *userstates.UserState) error {
	limit, err := s.toLimit(state)
	if err != nil {
		return err
	}
	return s.limitRepo.Save(ctx, limit)
}

func (s *saveLimitSaver) ClearData(state *userstates.UserState) {
	state.ClearBufferValue(userstates.SetLimitMonthIndex)
	state.ClearBufferValue(userstates.SetLimitMonthValue)
}
