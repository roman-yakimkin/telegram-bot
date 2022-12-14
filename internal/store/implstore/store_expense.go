package implstore

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/utils"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

func (s *store) calcAmount(ctx context.Context, data repo.ExpData, currTo string, conv convertors.CurrencyConvertorTo) (int, error) {
	sum := 0
	for _, currencyData := range data {
		for _, dayData := range currencyData {
			for date, payment := range dayData {
				converted, err := conv.To(ctx, payment, currTo, date)
				if err != nil {
					return 0, err
				}
				sum += converted
			}
		}
	}
	return sum, nil
}

func (s *store) amountPerYearMonth(ctx context.Context, userId int64, year int, month int, curr string, conv convertors.CurrencyConvertorTo) (int, error) {
	firstTime, lastTime := utils.FirstLastTimeOfMonth(year, month)
	expData, err := s.Expense().ExpensesByUserAndTimeInterval(ctx, userId, firstTime, lastTime)
	if err != nil {
		return 0, err
	}
	amount, err := s.calcAmount(ctx, expData, curr, conv)
	if err != nil {
		return 0, err
	}
	return amount, nil
}

func (s *store) MeetMonthlyLimit(ctx context.Context, userId int64, date time.Time, amountInRub int, conv repo.CurrencyConvertorTo) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "check monthly limit")
	defer span.Finish()

	y, m, _ := date.Date()
	amountAdded, err := s.amountPerYearMonth(ctx, userId, y, int(m), "RUB", conv)
	if err != nil {
		return false, err
	}
	limit, err := s.Limit().GetOne(ctx, userId, int(m))
	if err != nil {
		return false, err
	}
	if amountInRub+amountAdded > limit.Value {
		return false, nil
	}
	return true, nil
}
