package implstore

import (
	"context"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
)

func (s *store) UserCurrencyRate(ctx context.Context, UserID int64, date time.Time) (*currencies.CurrencyRate, error) {
	userInfo, err := s.UserState().GetOne(ctx, UserID)
	if err != nil {
		return nil, err
	}
	currency, err := s.CurrencyRate().GetOneByDate(ctx, userInfo.Currency, date)
	if err != nil {
		return nil, err
	}
	return currency, nil
}
