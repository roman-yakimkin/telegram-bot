package messages

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/reports"
)

const infoText = `/info - текущая справка
/newexpense - добавление новой траты
/lastweek - траты за последнюю неделю
/lastmonth - траты за последний месяц
/lastyear - траты за последний год`

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type Model struct {
	tgClient MessageSender
	rm       *reports.ReportManager
}

func New(tgClient MessageSender, rm *reports.ReportManager) *Model {
	return &Model{
		tgClient: tgClient,
		rm:       rm,
	}
}

type Message struct {
	Text   string
	UserID int64
}

func (s *Model) IncomingMessage(msg Message, status int) (int, error) {
	switch {
	case msg.Text == "/start":
		return userstates.ExpectedCommand, s.tgClient.SendMessage("hello\n"+infoText, msg.UserID)

	case msg.Text == "/info":
		return userstates.ExpectedCommand, s.tgClient.SendMessage(infoText, msg.UserID)

	case msg.Text == "/lastweek":
		return userstates.ExpectedCommand, s.tgClient.SendMessage(s.rm.LastWeek(msg.UserID), msg.UserID)

	case msg.Text == "/lastmonth":
		return userstates.ExpectedCommand, s.tgClient.SendMessage(s.rm.LastMonth(msg.UserID), msg.UserID)

	case msg.Text == "/lastyear":
		return userstates.ExpectedCommand, s.tgClient.SendMessage(s.rm.LastYear(msg.UserID), msg.UserID)

	case msg.Text == "/newexpense":
		return userstates.ExpectedCategory, s.tgClient.SendMessage("Введите категорию платежа", msg.UserID)

	case status == userstates.ExpectedCategory:
		return userstates.ExpectedAmount, s.tgClient.SendMessage("Введите сумму платежа", msg.UserID)

	case status == userstates.IncorrectCategory:
		return userstates.ExpectedCategory, s.tgClient.SendMessage("Категория задана неверно. Введите категорию платежа", msg.UserID)

	case status == userstates.ExpectedAmount:
		return userstates.ExpectedDate, s.tgClient.SendMessage("Введите дату платежа в формате ГГГГ-ММ-ДД (* - текущая дата)", msg.UserID)

	case status == userstates.IncorrectAmount:
		return userstates.ExpectedAmount, s.tgClient.SendMessage("Сумма платежа задана неверно. Введите сумму платежа", msg.UserID)

	case status == userstates.ExpectedDate:
		return userstates.ExpectedCommand, s.tgClient.SendMessage("Информация о платеже добавлена", msg.UserID)

	case status == userstates.IncorrectDate:
		return userstates.ExpectedDate, s.tgClient.SendMessage("Дата задана некорректно. Введите дату платежа в формате ГГГГ-ММ-ДД (* - текущая дата)", msg.UserID)

	default:
		return userstates.ExpectedCommand, s.tgClient.SendMessage("не знаю эту команду", msg.UserID)
	}
}
