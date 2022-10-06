package output

import (
	"fmt"
	"strings"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type CurrencyListOutput struct {
	currRepo repo.CurrencyRepo
}

func NewCurrencyOutput(currRepo repo.CurrencyRepo) *CurrencyListOutput {
	return &CurrencyListOutput{
		currRepo: currRepo,
	}
}

func (o *CurrencyListOutput) Output() (string, error) {
	var sb strings.Builder
	currencies, err := o.currRepo.GetAll()
	if err != nil {
		return "", nil
	}
	sb.WriteString("Введите валюту из списка или *, если передумали\n")
	for _, currency := range currencies {
		sb.WriteString(fmt.Sprintf("%s\n", currency.Name))
	}
	return sb.String(), nil
}
