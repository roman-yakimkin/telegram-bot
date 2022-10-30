package userstateprocessors

import (
	"context"
	"regexp"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/utils"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"go.uber.org/zap"
)

type setLimitAmountProcessor struct {
	processStatus int
	amountRegExp  *regexp.Regexp
	currConv      convertors.CurrencyConvertorFrom
	logger        *zap.Logger
}

func NewSetLimitAmountProcessor(currConv convertors.CurrencyConvertorFrom, logger *zap.Logger) (UserStateProcessor, error) {
	amountRegExp, err := regexp.Compile(AmountRegex)
	if err != nil {
		return nil, err
	}
	return &setLimitAmountProcessor{
		processStatus: userstates.ExpectedSetLimitAmount,
		amountRegExp:  amountRegExp,
		currConv:      currConv,
		logger:        logger,
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
		p.logger.Error("error on currency converting:", zap.Error(err))
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
