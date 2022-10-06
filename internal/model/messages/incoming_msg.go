package messages

import (
	"log"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

const infoText = `/info - текущая справка
/getcurrency - получение текущей валюты
/setcurrency - установка текущей валюты
/newexpense - добавление новой траты
/lastweek - траты за последнюю неделю
/lastmonth - траты за последний месяц
/lastyear - траты за последний год`

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type Model struct {
	tgClient MessageSender
	output   *output.Output
}

func New(tgClient MessageSender, output *output.Output) *Model {
	return &Model{
		tgClient: tgClient,
		output:   output,
	}
}

type Message struct {
	Text   string
	UserID int64
}

func (s *Model) reportMessage(UserID int64, msgText string) string {
	var report string
	var err error
	switch msgText {
	case "/lastweek":
		report, err = s.output.Reports().LastWeek(UserID)
	case "/lastmonth":
		report, err = s.output.Reports().LastMonth(UserID)
	case "/lastyear":
		report, err = s.output.Reports().LastYear(UserID)
	}
	if err != nil {
		log.Println("creating report error: ", err)
		report = "Ошибка при создании отчета"
	}
	return report
}

func (s *Model) IncomingMessage(msg Message, userState *userstates.UserState) (int, error) {
	switch {
	case msg.Text == "/start":
		return userstates.ExpectedCommand, s.tgClient.SendMessage("hello\n"+infoText, msg.UserID)

	case msg.Text == "/info":
		return userstates.ExpectedCommand, s.tgClient.SendMessage(infoText, msg.UserID)

	case msg.Text == "/getcurrency":
		return userstates.ExpectedCommand, s.tgClient.SendMessage("Ваша текущая валюта - "+userState.Currency, msg.UserID)

	case msg.Text == "/setcurrency":
		currOutput, err := s.output.Currency().Output()
		if err != nil {
			return userstates.ExpectedCommand, err
		}
		return userstates.ExpectedCurrency, s.tgClient.SendMessage(currOutput, msg.UserID)

	case userState.GetStatus() == userstates.ExpectedCurrency:
		return userstates.ExpectedCommand, s.tgClient.SendMessage("Валюта изменена", msg.UserID)

	case userState.GetStatus() == userstates.IncorrectCurrency:
		currOutput, err := s.output.Currency().Output()
		if err != nil {
			return userstates.ExpectedCommand, err
		}
		return userstates.ExpectedCurrency, s.tgClient.SendMessage("Валюта задана неверно\n"+currOutput, msg.UserID)

	case msg.Text == "/lastweek" || msg.Text == "/lastmonth" || msg.Text == "/lastyear":
		return userstates.ExpectedCommand, s.tgClient.SendMessage(s.reportMessage(msg.UserID, msg.Text), msg.UserID)

	case msg.Text == "/newexpense":
		return userstates.ExpectedCategory, s.tgClient.SendMessage("Введите категорию платежа", msg.UserID)

	case userState.GetStatus() == userstates.ExpectedCategory:
		return userstates.ExpectedAmount, s.tgClient.SendMessage("Введите сумму платежа. Текущая валюта - "+userState.Currency, msg.UserID)

	case userState.GetStatus() == userstates.IncorrectCategory:
		return userstates.ExpectedCategory, s.tgClient.SendMessage("Категория задана неверно. Введите категорию платежа", msg.UserID)

	case userState.GetStatus() == userstates.ExpectedAmount:
		return userstates.ExpectedDate, s.tgClient.SendMessage("Введите дату платежа в формате ГГГГ-ММ-ДД (* - текущая дата)", msg.UserID)

	case userState.GetStatus() == userstates.IncorrectAmount:
		return userstates.ExpectedAmount, s.tgClient.SendMessage("Сумма платежа задана неверно. Введите сумму платежа. Текущая валюта - "+userState.Currency, msg.UserID)

	case userState.GetStatus() == userstates.ExpectedDate:
		return userstates.ExpectedCommand, s.tgClient.SendMessage("Информация о платеже добавлена", msg.UserID)

	case userState.GetStatus() == userstates.IncorrectDate:
		return userstates.ExpectedDate, s.tgClient.SendMessage("Дата задана некорректно. Введите дату платежа в формате ГГГГ-ММ-ДД (* - текущая дата)", msg.UserID)

	default:
		return userstates.ExpectedCommand, s.tgClient.SendMessage("не знаю эту команду", msg.UserID)
	}
}
