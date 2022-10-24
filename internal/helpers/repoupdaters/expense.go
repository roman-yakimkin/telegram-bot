package repoupdaters

import (
	"context"
	"log"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type expenseSaver struct {
	expRepo      repo.ExpensesRepo
	limitChecker repo.ExpenseLimitChecker
}

func NewExpenseSaver(expRepo repo.ExpensesRepo, limitChecker repo.ExpenseLimitChecker) UserStateRepoUpdater {
	return &expenseSaver{
		expRepo:      expRepo,
		limitChecker: limitChecker,
	}
}

func (s *expenseSaver) toExpense(state *userstates.UserState) (*expenses.Expense, error) {
	amount, err := state.IfFloatTransformToInt(userstates.AddExpenseAmountValue)
	if err != nil {
		log.Println("error upon getting expense amount value", err)
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
