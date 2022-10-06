package implstore

import "gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"

func (s *Store) UserCurrency(UserID int64) (*currencies.Currency, error) {
	userInfo, err := s.UserState().GetOne(UserID)
	if err != nil {
		return nil, err
	}
	currency, err := s.Currency().GetOne(userInfo.Currency)
	if err != nil {
		return nil, err
	}
	return currency, nil
}
