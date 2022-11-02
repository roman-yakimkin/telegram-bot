package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

const (
	MessageNewExpense                   = "newexpense"
	MessageNewExpenseAmount             = "newexpense_amount"
	MessageNewExpenseCategory           = "newexpense_category"
	MessageNewExpenseCurrency           = "newexpense_currency"
	MessageNewExpenseDate               = "newexpense_date"
	MessageNewExpenseIncorrectAmount    = "newexpense_incorrectamount"
	MessageNewExpenseIncorrectCategory  = "newexpense_incorrectcategory"
	MessageNewExpenseIncorrectDate      = "newexpense_incorrectdate"
	MessageNewExpenseMonthLimitExceeded = "newexpense_monthlimitexceed"

	MessageDelLimit               = "dellimit"
	MessageDelLimitMonth          = "dellimit_month"
	MessageDelLimitIncorrectMonth = "dellimit_incorrectmonth"

	MessageSetLimit                = "setlimit"
	MessageSetLimitAmount          = "setlimit_amount"
	MessageSetLimitMonth           = "setlimit_month"
	MessageSetLimitIncorrectAmount = "setlimit_incorrectamount"
	MessageSetLimitIncorrectMonth  = "setlimit_incorrectmonth"

	MessageGetCurrency = "getcurrency"

	MessageSetCurrency                  = "setcurrency"
	MessageSetCurrencyIncorrectCurrency = "setcurrency_incorrectcurrency"

	MessageInfo   = "info"
	MessageLimits = "limits"
	MessageStart  = "start"
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
