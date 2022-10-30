package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type MessageSender interface {
	SendMessage(text string, userId int64) error
}

type Message struct {
	Text   string
	UserId int64
}

const InfoText = `/info - текущая справка
/getcurrency - получение текущей валюты
/setcurrency - установка текущей валюты
/setlimit - установить лимит за месяц
/dellimit - удалить лимит за месяц
/limits - получение лимитов по месяцам
/newexpense - добавление новой траты
/lastweek - траты за последнюю неделю
/lastmonth - траты за последний месяц
/lastyear - траты за последний год`

type MessageProcessor interface {
	ShouldProcess(msg Message, userState *userstates.UserState) bool
	DoProcess(ctx context.Context, msg Message, userState *userstates.UserState) (int, string, error)
}
