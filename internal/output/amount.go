package output

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type CurrencyAmount interface {
	Output(ctx context.Context, amount int, currName string) (string, error)
}

type currencyAmount struct {
	cr repo.CurrencyRepo
}

func NewCurrencyAmount(cr repo.CurrencyRepo) CurrencyAmount {
	return &currencyAmount{
		cr: cr,
	}
}

func (a *currencyAmount) Output(ctx context.Context, amount int, currName string) (string, error) {
	currency, err := a.cr.GetOne(ctx, currName)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf(currency.Display, float64(amount)/100)
	return result, nil
}
