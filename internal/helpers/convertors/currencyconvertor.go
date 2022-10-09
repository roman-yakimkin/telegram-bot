package convertors

import (
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

func (c *currencyConvertor) getRateToMain(currName string, date time.Time) (float64, error) {
	currency, err := c.currencyRepo.GetOneByDate(currName, date)
	if err != nil {
		return 0, err
	}
	return currency.RateToMain, nil
}

func (c *currencyConvertor) From(amount int, currFrom string, date time.Time) (int, error) {
	rate, err := c.getRateToMain(currFrom, date)
	if err != nil {
		return 0, err
	}
	return int(float64(amount) * rate), nil
}

func (c *currencyConvertor) To(amount int, currFrom string, date time.Time) (int, error) {
	rate, err := c.getRateToMain(currFrom, date)
	if err != nil {
		return 0, err
	}
	result := float64(amount) / rate
	return int(result), nil
}
