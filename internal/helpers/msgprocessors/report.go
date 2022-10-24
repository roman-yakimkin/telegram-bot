package msgprocessors

import (
	"context"
	"log"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type reportMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewReportMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &reportMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *reportMessageProcessor) ShouldProcess(msg Message, _ *userstates.UserState) bool {
	return msg.Text == "/lastweek" || msg.Text == "/lastmonth" || msg.Text == "/lastyear"
}

func (p *reportMessageProcessor) DoProcess(ctx context.Context, msg Message, _ *userstates.UserState) (int, error) {
	var report string
	var err error
	switch msg.Text {
	case "/lastweek":
		report, err = p.output.Reports().LastWeek(ctx, msg.UserId)
	case "/lastmonth":
		report, err = p.output.Reports().LastMonth(ctx, msg.UserId)
	case "/lastyear":
		report, err = p.output.Reports().LastYear(ctx, msg.UserId)
	}
	if err != nil {
		log.Println("creating report error: ", err)
		report = "Ошибка при создании отчета"
	}

	return userstates.ExpectedCommand, p.tgClient.SendMessage(report, msg.UserId)
}
