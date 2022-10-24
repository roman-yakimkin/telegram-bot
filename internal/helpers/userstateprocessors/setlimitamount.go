package userstateprocessors

import (
	"context"
	"log"
	"regexp"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/utils"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type setLimitAmountProcessor struct {
	processStatus int
	amountRegExp  *regexp.Regexp
	currConv      convertors.CurrencyConvertorFrom
}

func NewSetLimitAmountProcessor(currConv convertors.CurrencyConvertorFrom) (UserStateProcessor, error) {
	amountRegExp, err := regexp.Compile(AmountRegex)
	if err != nil {
		return nil, err
	}
	return &setLimitAmountProcessor{
		processStatus: userstates.ExpectedSetLimitAmount,
		amountRegExp:  amountRegExp,
		currConv:      currConv,
	}, nil
}

func (p *setLimitAmountProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *setLimitAmountProcessor) DoProcess(ctx context.Context, state *userstates.UserState, msgText string) {
	amountInt, amountFrac, err := utils.ParseMsgText(msgText, p.amountRegExp)
	if err != nil {
		state.SetStatus(userstates.IncorrectSetLimitAmount)
		return
	}
	amount := amountInt*100 + amountFrac
	if err = p.convertAndAddAmount(ctx, state, amount); err != nil {
		log.Println("error on currency converting:", err)
		state.SetStatus(userstates.ExpectedCommand)
		return
	}
}

func (p *setLimitAmountProcessor) convertAndAddAmount(ctx context.Context, state *userstates.UserState, unconvertedAmount int) error {
	amountInBaseCurrency, err := p.currConv.From(ctx, unconvertedAmount, state.Currency, utils.TimeTruncate(time.Now()))
	if err != nil {
		return err
	}
	state.SetBufferValue(userstates.SetLimitMonthValue, amountInBaseCurrency)
	return nil
}
