package msgprocessors

import "gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type Message struct {
	Text   string
	UserID int64
}

const InfoText = `/info - текущая справка
/getcurrency - получение текущей валюты
/setcurrency - установка текущей валюты
/newexpense - добавление новой траты
/lastweek - траты за последнюю неделю
/lastmonth - траты за последний месяц
/lastyear - траты за последний год`

type MessageProcessor interface {
	ShouldProcess(msg Message, userState *userstates.UserState) bool
	DoProcess(msg Message, userState *userstates.UserState) (int, error)
}
