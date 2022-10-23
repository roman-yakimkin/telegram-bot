package userstateprocessors

import (
	"context"
	"strings"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type currencyProcessor struct {
	processStatus int
	currRepo      repo.CurrencyRepo
}

func NewCurrencyProcessor(currRepo repo.CurrencyRepo) UserStateProcessor {
	return &currencyProcessor{
		processStatus: userstates.ExpectedCurrency,
		currRepo:      currRepo,
	}
}

func (p *currencyProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *currencyProcessor) DoProcess(ctx context.Context, state *userstates.UserState, msgText string) {
	if msgText == "*" {
		state.SetStatus(userstates.ExpectedCommand)
		return
	}
	currencies, err := p.currRepo.GetAll(ctx)
	if err != nil {
		state.SetStatus(userstates.IncorrectCurrency)
		return
	}
	for _, currency := range currencies {
		if strings.EqualFold(currency.Name, msgText) {
			state.Currency = strings.ToUpper(currency.Name)
			return
		}
	}
	state.SetStatus(userstates.IncorrectCurrency)
}
