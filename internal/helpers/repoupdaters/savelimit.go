package repoupdaters

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
	"go.uber.org/zap"
)

type saveLimitSaver struct {
	limitRepo repo.ExpenseLimitsRepo
	logger    *zap.Logger
}

func NewSaveLimitSaver(limitRepo repo.ExpenseLimitsRepo, logger *zap.Logger) UserStateRepoUpdater {
	return &saveLimitSaver{
		limitRepo: limitRepo,
		logger:    logger,
	}
}

func (s *saveLimitSaver) toLimit(state *userstates.UserState) (*expenses.ExpenseLimit, error) {
	month, err := state.IfFloatTransformToInt(userstates.SetLimitMonthIndex)
	if err != nil {
		s.logger.Error("error upon getting month index", zap.Error(err))
		return nil, err
	}
	value, err := state.IfFloatTransformToInt(userstates.SetLimitMonthValue)
	if err != nil {
		s.logger.Error("error upon getting month value", zap.Error(err))
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "update limit")
	defer span.Finish()

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
