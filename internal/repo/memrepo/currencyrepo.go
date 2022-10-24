package memrepo

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/currencies"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type currencyRepo struct {
	c map[string]currencies.Currency
}

func NewCurrencyRepo(service *config.Service) (repo.CurrencyRepo, error) {
	cfg := service.GetConfig()
	currMap := make(map[string]currencies.Currency, len(cfg.Currencies))
	for _, curr := range cfg.Currencies {
		currMap[curr.Name] = currencies.Currency{
			Name:    curr.Name,
			Display: curr.Display,
		}
	}
	return &currencyRepo{
		c: currMap,
	}, nil
}

func (r *currencyRepo) GetOne(_ context.Context, currName string) (*currencies.Currency, error) {
	curr, ok := r.c[currName]
	if !ok {
		return nil, localerr.ErrCurrencyNotFound
	}
	return &curr, nil
}

func (r *currencyRepo) GetAll(_ context.Context) ([]currencies.Currency, error) {
	result := make([]currencies.Currency, 0, len(r.c))
	for _, curr := range r.c {
		result = append(result, curr)
	}
	return result, nil
}
