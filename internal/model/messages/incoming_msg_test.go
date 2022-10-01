package messages

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockmessages "gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/mocks/messages"
)

func Test_OnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mockmessages.NewMockMessageSender(ctrl)
	model := New(sender, nil)

	sender.EXPECT().SendMessage("hello\n"+infoText, int64(123))

	_, err := model.IncomingMessage(Message{
		Text:   "/start",
		UserID: 123,
	}, userstates.ExpectedCommand)

	assert.NoError(t, err)
}

func Test_OnInfoCommand_ShouldAnswerWithInfoMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mockmessages.NewMockMessageSender(ctrl)
	model := New(sender, nil)

	sender.EXPECT().SendMessage(infoText, int64(123))

	_, err := model.IncomingMessage(Message{
		Text:   "/info",
		UserID: 123,
	}, userstates.ExpectedCommand)

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
	}, userstates.ExpectedCommand)

	assert.NoError(t, err)
}
