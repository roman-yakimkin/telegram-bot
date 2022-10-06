package messages

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockmessages "gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/mocks/messages"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

func Test_OnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mockmessages.NewMockMessageSender(ctrl)
	model := New(sender, nil)
	var uid int64 = 123
	userState := userstates.NewUserState(uid)

	sender.EXPECT().SendMessage("hello\n"+infoText, uid)

	_, err := model.IncomingMessage(Message{
		Text:   "/start",
		UserID: uid,
	}, userState)

	assert.NoError(t, err)
}

func Test_OnInfoCommand_ShouldAnswerWithInfoMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mockmessages.NewMockMessageSender(ctrl)
	model := New(sender, nil)
	var uid int64 = 123
	userState := userstates.NewUserState(uid)

	sender.EXPECT().SendMessage(infoText, uid)

	_, err := model.IncomingMessage(Message{
		Text:   "/info",
		UserID: uid,
	}, userState)

	assert.NoError(t, err)
}

func Test_OnUnknownCommand_ShouldAnswerWithHelpMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mockmessages.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage("не знаю эту команду", int64(123))
	model := New(sender, nil)
	var uid int64 = 123
	userState := userstates.NewUserState(uid)

	_, err := model.IncomingMessage(Message{
		Text:   "some text",
		UserID: uid,
	}, userState)

	assert.NoError(t, err)
}
