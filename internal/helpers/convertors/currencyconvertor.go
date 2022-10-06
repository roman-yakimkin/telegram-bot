package convertors

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type Currency struct {
	cr  repo.CurrencyRepo
	cfg *config.Service
}

func NewCurrency(cr repo.CurrencyRepo, cfg *config.Service) *Currency {
	return &Currency{
		cr:  cr,
		cfg: cfg,
	}
}

func (c *Currency) getRateToMain(currName string) (float64, error) {
	currency, err := c.cr.GetOne(currName)
	if err != nil {
		return 0, err
	}
	return currency.RateToMain, nil
}

func (c *Currency) From(amount int, currFrom string) (int, error) {
	rate, err := c.getRateToMain(currFrom)
	if err != nil {
		return 0, err
	}
	return int(float64(amount) * rate), nil
}

func (c *Currency) To(amount int, currFrom string) (int, error) {
	rate, err := c.getRateToMain(currFrom)
	if err != nil {
		return 0, err
	}
	result := float64(amount) / rate
	return int(result), nil
}
