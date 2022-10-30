package repoupdaters

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
	"go.uber.org/zap"
)

type expenseSaver struct {
	expRepo      repo.ExpensesRepo
	limitChecker repo.ExpenseLimitChecker
	logger       *zap.Logger
}

func NewExpenseSaver(expRepo repo.ExpensesRepo, limitChecker repo.ExpenseLimitChecker, logger *zap.Logger) UserStateRepoUpdater {
	return &expenseSaver{
		expRepo:      expRepo,
		limitChecker: limitChecker,
		logger:       logger,
	}
}

func (s *expenseSaver) toExpense(state *userstates.UserState) (*expenses.Expense, error) {
	amount, err := state.IfFloatTransformToInt(userstates.AddExpenseAmountValue)
	if err != nil {
		s.logger.Error("error upon getting expense amount value", zap.Error(err))
		return nil, err
	}
	return &expenses.Expense{
		UserId:   state.UserId,
		Category: state.GetBufferValue(userstates.AddExpenseCategoryValue).(string),
		Amount:   amount,
		Currency: state.Currency,
		Date:     state.GetBufferValue(userstates.AddExpenseDateValue).(time.Time),
	}, nil
}

func (s *expenseSaver) ReadyToUpdate(state *userstates.UserState) bool {
	ok1 := state.BufferValueExists(userstates.AddExpenseCategoryValue)
	ok2 := state.BufferValueExists(userstates.AddExpenseAmountValue)
	ok3 := state.BufferValueExists(userstates.AddExpenseDateValue)
	return ok1 && ok2 && ok3
}

func (s *expenseSaver) UpdateRepo(ctx context.Context, state *userstates.UserState) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "save expense")
	defer span.Finish()

	expense, err := s.toExpense(state)
	if err != nil {
		return err
	}
	return s.expRepo.Add(ctx, expense, s.limitChecker)
}

func (s *expenseSaver) ClearData(state *userstates.UserState) {
	state.ClearBufferValue(userstates.AddExpenseCategoryValue)
	state.ClearBufferValue(userstates.AddExpenseAmountValue)
	state.ClearBufferValue(userstates.AddExpenseDateValue)
}
