package repoupdaters

import (
	"context"
	"log"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type delLimitSaver struct {
	limitRepo repo.ExpenseLimitsRepo
}

func NewDelLimitSaver(limitRepo repo.ExpenseLimitsRepo) UserStateRepoUpdater {
	return &delLimitSaver{
		limitRepo: limitRepo,
	}
}

func (s *delLimitSaver) toLimitMonth(state *userstates.UserState) (int, error) {
	index, err := state.IfFloatTransformToInt(userstates.DeleteLimitMonthIndex)
	if err != nil {
		log.Println("error upon getting limit month index: ", err)
	}
	return index, err
}

func (s *delLimitSaver) ReadyToUpdate(state *userstates.UserState) bool {
	return state.BufferValueExists(userstates.DeleteLimitMonthIndex)
}

func (s *delLimitSaver) UpdateRepo(ctx context.Context, state *userstates.UserState) error {
	index, err := s.toLimitMonth(state)
	if err != nil {
		return err
	}
	return s.limitRepo.Delete(ctx, state.UserID, index)
}

func (s *delLimitSaver) ClearData(state *userstates.UserState) {
	state.ClearBufferValue(userstates.DeleteLimitMonthIndex)
}
