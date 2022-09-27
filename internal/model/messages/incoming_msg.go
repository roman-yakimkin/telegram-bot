package messages

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/reports"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/vars"
)

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
		return vars.ExpectedCommand, s.tgClient.SendMessage("hello", msg.UserID)

	case msg.Text == "/info":
		text := `/info - текущая справка
/newexpense - добавление новой траты
/lastweek - траты за последнюю неделю
/lastmonth - траты за последний месяц
/lastyear - траты за последний год`
		return vars.ExpectedCommand, s.tgClient.SendMessage(text, msg.UserID)

	case msg.Text == "/lastweek":
		return vars.ExpectedCommand, s.tgClient.SendMessage(s.rm.LastWeek(msg.UserID), msg.UserID)

	case msg.Text == "/lastmonth":
		return vars.ExpectedCommand, s.tgClient.SendMessage(s.rm.LastMonth(msg.UserID), msg.UserID)

	case msg.Text == "/lastyear":
		return vars.ExpectedCommand, s.tgClient.SendMessage(s.rm.LastYear(msg.UserID), msg.UserID)

	case msg.Text == "/newexpense":
		return vars.ExpectedCategory, s.tgClient.SendMessage("Введите категорию платежа", msg.UserID)

	case status == vars.ExpectedCategory:
		return vars.ExpectedAmount, s.tgClient.SendMessage("Введите сумму платежа", msg.UserID)

	case status == vars.IncorrectCategory:
		return vars.ExpectedCategory, s.tgClient.SendMessage("Категория задана неверно. Введите категорию платежа", msg.UserID)

	case status == vars.ExpectedAmount:
		return vars.ExpectedDate, s.tgClient.SendMessage("Введите дату платежа в формате ГГГГ-ММ-ДД (* - текущая дата)", msg.UserID)

	case status == vars.IncorrectAmount:
		return vars.ExpectedAmount, s.tgClient.SendMessage("Сумма платежа задана неверно. Введите сумму платежа", msg.UserID)

	case status == vars.ExpectedDate:
		return vars.ExpectedCommand, s.tgClient.SendMessage("Информация о платеже добавлена", msg.UserID)

	case status == vars.IncorrectDate:
		return vars.ExpectedDate, s.tgClient.SendMessage("Дата задана некорректно. Введите дату платежа в формате ГГГГ-ММ-ДД (* - текущая дата)", msg.UserID)

	default:
		return vars.ExpectedCommand, s.tgClient.SendMessage("не знаю эту команду", msg.UserID)
	}
}
