package output

import (
	"context"
	"fmt"
	"strings"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type CurrencyListOutput interface {
	Output(ctx context.Context) (string, error)
}

type currencyListOutput struct {
	currRepo repo.CurrencyRepo
}

func NewCurrencyListOutput(currRepo repo.CurrencyRepo) CurrencyListOutput {
	return &currencyListOutput{
		currRepo: currRepo,
	}
}

func (o *currencyListOutput) Output(ctx context.Context) (string, error) {
	var sb strings.Builder
	currencies, err := o.currRepo.GetAll(ctx)
	if err != nil {
		return "", nil
	}
	sb.WriteString("Введите валюту из списка или *, если передумали\n")
	for _, currency := range currencies {
		sb.WriteString(fmt.Sprintf("%s\n", currency.Name))
	}
	return sb.String(), nil
}
