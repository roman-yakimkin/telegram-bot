package output

import (
	"bytes"
	"fmt"
	"text/template"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type CurrencyAmount struct {
	cr repo.CurrencyRepo
}

func NewCurrencyAmount(cr repo.CurrencyRepo) *CurrencyAmount {
	return &CurrencyAmount{
		cr: cr,
	}
}

func (a *CurrencyAmount) Output(amount int, currName string) (string, error) {
	currency, err := a.cr.GetOne(currName)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	t := template.Must(template.New("display").Parse(currency.Display))
	params := map[string]interface{}{
		"amount": fmt.Sprintf("%.2f", float64(amount)/100),
	}
	err = t.Execute(&b, params)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
