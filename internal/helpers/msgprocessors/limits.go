package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
	"go.uber.org/zap"
)

type limitsMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
	logger   *zap.Logger
}

func NewLimitsMessageProcessor(ms MessageSender, output output.Output, logger *zap.Logger) MessageProcessor {
	return &limitsMessageProcessor{
		tgClient: ms,
		output:   output,
		logger:   logger,
	}
}

func (p *limitsMessageProcessor) ShouldProcess(msg Message, _ *userstates.UserState) bool {
	return msg.Text == "/limits"
}

func (p *limitsMessageProcessor) DoProcess(ctx context.Context, msg Message, _ *userstates.UserState) (int, string, error) {
	limits, err := p.output.Limits().Output(ctx, msg.UserId)
	if err != nil {
		limits = "Ошибка при выводе лимитов"
		p.logger.Error("limit output error", zap.Error(err))
	}
	return userstates.ExpectedCommand, MessageLimits, p.tgClient.SendMessage(limits, msg.UserId)
}
