package repoupdaters

import (
	"log"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type expenseSaver struct {
	userState *userstates.UserState
	expRepo   repo.ExpensesRepo
}

func NewExpenseSaver(expRepo repo.ExpensesRepo) UserStateRepoUpdater {
	return &expenseSaver{
		expRepo: expRepo,
	}
}

func (s *expenseSaver) SetUserState(userState *userstates.UserState) {
	s.userState = userState
}

func (s *expenseSaver) toExpense() (*expenses.Expense, error) {
	amount, err := s.userState.IfFloatTransformToInt(userstates.AddExpenseAmountValue)
	if err != nil {
		log.Println("error upon getting expense amount value", err)
		return nil, err
	}
	return &expenses.Expense{
		UserID:   s.userState.UserID,
		Category: s.userState.GetBufferValue(userstates.AddExpenseCategoryValue).(string),
		Amount:   amount,
		Currency: s.userState.Currency,
		Date:     s.userState.GetBufferValue(userstates.AddExpenseDateValue).(time.Time),
	}, nil
}

func (s *expenseSaver) ReadyToUpdate() bool {
	ok1 := s.userState.BufferValueExists(userstates.AddExpenseCategoryValue)
	ok2 := s.userState.BufferValueExists(userstates.AddExpenseAmountValue)
	ok3 := s.userState.BufferValueExists(userstates.AddExpenseDateValue)
	return ok1 && ok2 && ok3
}

func (s *expenseSaver) UpdateRepo() error {
	expense, err := s.toExpense()
	if err != nil {
		return err
	}
	return s.expRepo.Add(expense)
}

func (s *expenseSaver) ClearData() {
	s.userState.ClearBufferValue(userstates.AddExpenseCategoryValue)
	s.userState.ClearBufferValue(userstates.AddExpenseAmountValue)
	s.userState.ClearBufferValue(userstates.AddExpenseDateValue)
}
