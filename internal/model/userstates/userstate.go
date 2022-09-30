package userstates

import (
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
)

const (
	ExpectedCommand = iota
	ExpectedCategory
	IncorrectCategory
	ExpectedAmount
	IncorrectAmount
	ExpectedDate
	IncorrectDate
)

type UserState struct {
	UserID        int64
	status        int
	category      string
	amount        int
	date          time.Time
	addedCategory bool
	addedAmount   bool
	addedDate     bool
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

func (s *UserState) GetStatus() int {
	return s.status
}

func (s *UserState) SetStatus(newStatus int) {
	s.status = newStatus
}

func (s *UserState) Added() bool {
	return s.addedCategory && s.addedAmount && s.addedDate
}

func (s *UserState) ToExpense() *expenses.Expense {
	return &expenses.Expense{
		UserID:   s.UserID,
		Category: s.category,
		Amount:   s.amount,
		Date:     s.date,
	}
}
