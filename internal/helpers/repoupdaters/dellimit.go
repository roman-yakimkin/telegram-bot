package repoupdaters

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
	"go.uber.org/zap"
)

type delLimitSaver struct {
	limitRepo repo.ExpenseLimitsRepo
	logger    *zap.Logger
}

func NewDelLimitSaver(limitRepo repo.ExpenseLimitsRepo, logger *zap.Logger) UserStateRepoUpdater {
	return &delLimitSaver{
		limitRepo: limitRepo,
		logger:    logger,
	}
}

func (s *delLimitSaver) toLimitMonth(state *userstates.UserState) (int, error) {
	index, err := state.IfFloatTransformToInt(userstates.DeleteLimitMonthIndex)
	if err != nil {
		s.logger.Error("error upon getting limit month index: ", zap.Error(err))
	}
	return index, err
}

func (s *delLimitSaver) ReadyToUpdate(state *userstates.UserState) bool {
	return state.BufferValueExists(userstates.DeleteLimitMonthIndex)
}

func (s *delLimitSaver) UpdateRepo(ctx context.Context, state *userstates.UserState) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "delete limit")
	defer span.Finish()

	index, err := s.toLimitMonth(state)
	if err != nil {
		return err
	}
	return s.limitRepo.Delete(ctx, state.UserId, index)
}

func (s *delLimitSaver) ClearData(state *userstates.UserState) {
	state.ClearBufferValue(userstates.DeleteLimitMonthIndex)
}
