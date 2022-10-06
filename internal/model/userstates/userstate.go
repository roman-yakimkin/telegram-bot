package userstates

import (
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
)

const (
	ExpectedCommand = iota
	ExpectedCurrency
	IncorrectCurrency
	ExpectedCategory
	IncorrectCategory
	ExpectedAmount
	IncorrectAmount
	ExpectedDate
	IncorrectDate
)

type UserState struct {
	UserID        int64
	Currency      string
	status        int
	category      string
	amount        int
	date          time.Time
	addedCategory bool
	addedAmount   bool
	addedDate     bool
}

func NewUserState(UserID int64) *UserState {
	return &UserState{
		UserID:   UserID,
		Currency: "RUB",
		status:   ExpectedCommand,
	}
}

func (s *UserState) BeforeSave() {
	if s.Currency == "" {
		s.Currency = "RUB"
	}
}

func (s *UserState) GetStatus() int {
	return s.status
}

func (s *UserState) cleanInputtedExpense() {
	s.category = ""
	s.amount = 0
	s.date = time.Time{}
	s.addedCategory = false
	s.addedAmount = false
	s.addedDate = false
}

func (s *UserState) SetStatus(newStatus int) {
	s.status = newStatus
	switch newStatus {
	case ExpectedCommand:
		s.cleanInputtedExpense()
	}
}

func (s *UserState) SetCategory(category string) {
	s.category = category
	s.addedCategory = true
}

func (s *UserState) SetAmount(amount int) {
	s.amount = amount
	s.addedAmount = true
}

func (s *UserState) SetDate(date time.Time) {
	s.date = date
	s.addedDate = true
}

func (s *UserState) Added() bool {
	return s.addedCategory && s.addedAmount && s.addedDate
}

func (s *UserState) ToExpense() *expenses.Expense {
	return &expenses.Expense{
		UserID:   s.UserID,
		Category: s.category,
		Amount:   s.amount,
		Currency: s.Currency,
		Date:     s.date,
	}
}
