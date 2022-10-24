package convertors

import (
	"context"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type currencyConvertor struct {
	currencyRepo repo.CurrencyRateRepo
	cfg          *config.Service
}

func NewCurrencyConvertor(currencyRepo repo.CurrencyRateRepo, cfg *config.Service) CurrencyConvertor {
	return &currencyConvertor{
		currencyRepo: currencyRepo,
		cfg:          cfg,
	}
}

func (c *currencyConvertor) getRateToMain(ctx context.Context, currName string, date time.Time) (float64, error) {
	currency, err := c.currencyRepo.GetOneByDate(ctx, currName, date)
	if err != nil {
		return 0, err
	}
	return currency.RateToMain, nil
}

func (c *currencyConvertor) From(ctx context.Context, amount int, currFrom string, date time.Time) (int, error) {
	rate, err := c.getRateToMain(ctx, currFrom, date)
	if err != nil {
		return 0, err
	}
	return int(float64(amount) * rate), nil
}

func (c *currencyConvertor) To(ctx context.Context, amount int, currFrom string, date time.Time) (int, error) {
	rate, err := c.getRateToMain(ctx, currFrom, date)
	if err != nil {
		return 0, err
	}
	result := float64(amount) / rate
	return int(result), nil
}
