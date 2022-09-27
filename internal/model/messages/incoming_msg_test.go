package messages

import (
	mockmessages "gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/mocks/messages"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/vars"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_OnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mockmessages.NewMockMessageSender(ctrl)
	model := New(sender, nil)

	sender.EXPECT().SendMessage("hello", int64(123))

	_, err := model.IncomingMessage(Message{
		Text:   "/start",
		UserID: 123,
	}, vars.ExpectedCommand)

	assert.NoError(t, err)
}

func Test_OnInfoCommand_ShouldAnswerWithInfoMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mockmessages.NewMockMessageSender(ctrl)
	model := New(sender, nil)

	text := `/info - текущая справка
/newexpense - добавление новой траты
/lastweek - траты за последнюю неделю
/lastmonth - траты за последний месяц
/lastyear - траты за последний год`

	sender.EXPECT().SendMessage(text, int64(123))

	_, err := model.IncomingMessage(Message{
		Text:   "/info",
		UserID: 123,
	}, vars.ExpectedCommand)

	assert.NoError(t, err)
}

func Test_OnUnknownCommand_ShouldAnswerWithHelpMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mockmessages.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage("не знаю эту команду", int64(123))
	model := New(sender, nil)

	_, err := model.IncomingMessage(Message{
		Text:   "some text",
		UserID: 123,
	}, vars.ExpectedCommand)

	assert.NoError(t, err)
}
