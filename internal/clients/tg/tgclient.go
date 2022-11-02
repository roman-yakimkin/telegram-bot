package tg

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/localmetrics"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/msgprocessors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/repoupdaters"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/userstateprocessors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/store"
	"go.uber.org/zap"
)

type TokenGetter interface {
	Token() string
}

type Client struct {
	client   *tgbotapi.BotAPI
	store    store.Store
	currConv convertors.CurrencyConvertor
	logger   *zap.Logger
}

func New(tokenGetter TokenGetter, store store.Store, currConv convertors.CurrencyConvertor, logger *zap.Logger) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	return &Client{
		client:   client,
		store:    store,
		currConv: currConv,
		logger:   logger,
	}, nil
}

func (c *Client) SendMessage(text string, userId int64) error {
	_, err := c.client.Send(tgbotapi.NewMessage(userId, text))
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}
	return nil
}

func (c *Client) ListenUpdates(ctx context.Context, msgModel *messages.Model) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	amountProcessor, err := userstateprocessors.NewAmountProcessor()
	if err != nil {
		return err
	}
	setLimitAmountProcessor, err := userstateprocessors.NewSetLimitAmountProcessor(c.currConv, c.logger)
	if err != nil {
		return err
	}
	userStateProcessors := []userstateprocessors.UserStateProcessor{
		userstateprocessors.NewCategoryProcessor(),
		amountProcessor,
		userstateprocessors.NewDateProcessor(c.store, c.currConv, c.logger),
		userstateprocessors.NewCurrencyProcessor(c.store.Currency()),
		userstateprocessors.NewSetLimitMonthProcessor(),
		setLimitAmountProcessor,
		userstateprocessors.NewDelLimitMonthProcessor(),
	}
	repoUpdaters := []repoupdaters.UserStateRepoUpdater{
		repoupdaters.NewExpenseSaver(c.store.Expense(), c.store, c.logger),
		repoupdaters.NewSaveLimitSaver(c.store.Limit(), c.logger),
		repoupdaters.NewDelLimitSaver(c.store.Limit(), c.logger),
	}

	updates := c.client.GetUpdatesChan(u)
	c.logger.Info("Listen for updates")

	for update := range updates {
		if update.Message != nil { // If we got a message
			c.logger.Info(fmt.Sprintf("[%s] %s", update.Message.From.UserName, update.Message.Text))
			startTime := time.Now()
			span, ctx := opentracing.StartSpanFromContext(ctx, "undefined message")

			uid := update.Message.From.ID
			text := update.Message.Text
			userState, err := c.store.UserState().GetOne(ctx, uid)
			if err != nil && errors.Is(err, localerr.ErrUserStateNotFound) {
				userState = userstates.NewUserState(uid)
			}
			if err == nil && userState.GetStatus() != userstates.ExpectedCommand {
				for _, proc := range userStateProcessors {
					if userState.GetStatus() == proc.GetProcessStatus() {
						span, ctx := opentracing.StartSpanFromContext(ctx, "performing state processor")
						proc.DoProcess(ctx, userState, text)
						span.Finish()
						break
					}
				}
				for _, updater := range repoUpdaters {
					if updater.ReadyToUpdate(userState) {
						err := updater.UpdateRepo(ctx, userState)
						if err != nil {
							c.logger.Error("repo update error:", zap.Error(err))
						}
						updater.ClearData(userState)
					}
				}
			}

			newStatus, messageId, err := msgModel.IncomingMessage(ctx, msgprocessors.Message{
				Text:   text,
				UserId: uid,
			}, userState)

			span.SetOperationName(messageId)

			localmetrics.CntMessages.WithLabelValues(fmt.Sprintf("%d", uid), messageId).Inc()

			if err != nil {
				c.logger.Error("error processing message:", zap.Error(err))
			}
			userState.SetStatus(newStatus)
			err = c.store.UserState().Save(ctx, userState)
			if err != nil {
				c.logger.Error("error saving user state:", zap.Error(err))
			}

			span.Finish()

			diff := time.Since(startTime).Microseconds()
			localmetrics.PerformDuration.WithLabelValues(fmt.Sprintf("%d", uid), messageId).Observe(float64(diff))
		}
	}
	return nil
}
