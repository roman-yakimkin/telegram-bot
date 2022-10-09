package msgprocessors

import (
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

func (p *reportMessageProcessor) DoProcess(msg Message, _ *userstates.UserState) (int, error) {
	var report string
	var err error
	switch msg.Text {
	case "/lastweek":
		report, err = p.output.Reports().LastWeek(msg.UserID)
	case "/lastmonth":
		report, err = p.output.Reports().LastMonth(msg.UserID)
	case "/lastyear":
		report, err = p.output.Reports().LastYear(msg.UserID)
	}
	if err != nil {
		log.Println("creating report error: ", err)
		report = "Ошибка при создании отчета"
	}

	return userstates.ExpectedCommand, p.tgClient.SendMessage(report, msg.UserID)
}
