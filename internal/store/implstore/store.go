package implstore

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type Store struct {
	er repo.ExpensesRepo
	us repo.UserStateRepo
	cr repo.CurrencyRepo
}

func NewStore(er repo.ExpensesRepo, us repo.UserStateRepo, cr repo.CurrencyRepo) *Store {
	return &Store{
		er: er,
		us: us,
		cr: cr,
	}
}

func (s *Store) Expense() repo.ExpensesRepo {
	return s.er
}

func (s *Store) UserState() repo.UserStateRepo {
	return s.us
}

func (s *Store) Currency() repo.CurrencyRepo {
	return s.cr
}
