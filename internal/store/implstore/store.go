package implstore

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
	pkgstore "gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/store"
)

type store struct {
	er repo.ExpensesRepo
	us repo.UserStateRepo
	cu repo.CurrencyRepo
	cr repo.CurrencyRateRepo
	el repo.ExpenseLimitsRepo
}

func NewStore(er repo.ExpensesRepo, us repo.UserStateRepo, cu repo.CurrencyRepo, cr repo.CurrencyRateRepo, el repo.ExpenseLimitsRepo) pkgstore.Store {
	return &store{
		er: er,
		us: us,
		cu: cu,
		cr: cr,
		el: el,
	}
}

func (s *store) Expense() repo.ExpensesRepo {
	return s.er
}

func (s *store) UserState() repo.UserStateRepo {
	return s.us
}

func (s *store) Currency() repo.CurrencyRepo {
	return s.cu
}

func (s *store) CurrencyRate() repo.CurrencyRateRepo {
	return s.cr
}

func (s *store) Limit() repo.ExpenseLimitsRepo {
	return s.el
}
