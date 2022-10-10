package implstore

import (
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
)

func (s *store) UserCurrencyRate(UserID int64, date time.Time) (*currencies.CurrencyRate, error) {
	userInfo, err := s.UserState().GetOne(UserID)
	if err != nil {
		return nil, err
	}
	currency, err := s.CurrencyRate().GetOneByDate(userInfo.Currency, date)
	if err != nil {
		return nil, err
	}
	return currency, nil
}
