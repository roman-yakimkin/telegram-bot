package messages

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/msgprocessors"
	mockmessages "gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/mocks/messages"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"go.uber.org/zap"
)

func Test_OnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	sender := mockmessages.NewMockMessageSender(ctrl)
	logger, _ := zap.NewDevelopment()

	model := New(sender, nil, logger)
	var uid int64 = 123
	userState := userstates.NewUserState(uid)

	sender.EXPECT().SendMessage("hello\n"+msgprocessors.InfoText, uid)

	_, _, err := model.IncomingMessage(ctx, msgprocessors.Message{
		Text:   "/start",
		UserId: uid,
	}, userState)

	assert.NoError(t, err)
}

func Test_OnInfoCommand_ShouldAnswerWithInfoMessage(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	logger, _ := zap.NewDevelopment()

	sender := mockmessages.NewMockMessageSender(ctrl)
	model := New(sender, nil, logger)
	var uid int64 = 123
	userState := userstates.NewUserState(uid)

	sender.EXPECT().SendMessage(msgprocessors.InfoText, uid)

	_, _, err := model.IncomingMessage(ctx, msgprocessors.Message{
		Text:   "/info",
		UserId: uid,
	}, userState)

	assert.NoError(t, err)
}

func Test_OnUnknownCommand_ShouldAnswerWithHelpMessage(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	logger, _ := zap.NewDevelopment()

	sender := mockmessages.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage("не знаю эту команду", int64(123))
	model := New(sender, nil, logger)
	var uid int64 = 123
	userState := userstates.NewUserState(uid)

	_, _, err := model.IncomingMessage(ctx, msgprocessors.Message{
		Text:   "some text",
		UserId: uid,
	}, userState)

	assert.NoError(t, err)
}
