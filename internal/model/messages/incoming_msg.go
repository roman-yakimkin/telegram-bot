package messages

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/msgprocessors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type Model struct {
	tgClient      msgprocessors.MessageSender
	output        output.Output
	msgProcessors []msgprocessors.MessageProcessor
}

func New(tgClient msgprocessors.MessageSender, output output.Output) *Model {
	msgProcessors := []msgprocessors.MessageProcessor{
		msgprocessors.NewStartMessageProcessor(tgClient, output),
		msgprocessors.NewInfoMessageProcessor(tgClient, output),
		msgprocessors.NewGetCurrencyMessageProcessor(tgClient, output),
		msgprocessors.NewSetCurrencyMessageProcessor(tgClient, output),
		msgprocessors.NewExpectedCurrencyMessageProcessor(tgClient, output),
		msgprocessors.NewIncorrectCurrencyMessageProcessor(tgClient, output),
		msgprocessors.NewReportMessageProcessor(tgClient, output),
		msgprocessors.NewNewExpenseMessageProcessor(tgClient, output),
		msgprocessors.NewExpectedCategoryMessageProcessor(tgClient, output),
		msgprocessors.NewIncorrectCategoryMessageProcessor(tgClient, output),
		msgprocessors.NewExpectedAmountMessageProcessor(tgClient, output),
		msgprocessors.NewIncorrectAmountMessageProcessor(tgClient, output),
		msgprocessors.NewExpectedDateMessageProcessor(tgClient, output),
		msgprocessors.NewIncorrectDateMessageProcessor(tgClient, output),
	}
	return &Model{
		tgClient:      tgClient,
		output:        output,
		msgProcessors: msgProcessors,
	}
}

func (s *Model) IncomingMessage(msg msgprocessors.Message, userState *userstates.UserState) (int, error) {
	for _, proc := range s.msgProcessors {
		if proc.ShouldProcess(msg, userState) {
			return proc.DoProcess(msg, userState)
		}
	}
	return userstates.ExpectedCommand, s.tgClient.SendMessage("не знаю эту команду", msg.UserID)
}
