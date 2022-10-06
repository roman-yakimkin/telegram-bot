package tg

import (
	"github.com/pkg/errors"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/userstateprocessors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/store"
)

type TokenGetter interface {
	Token() string
}

type Client struct {
	client   *tgbotapi.BotAPI
	store    store.Store
	currConv convertors.CurrencyConvertorFrom
}

func New(tokenGetter TokenGetter, store store.Store, currConv convertors.CurrencyConvertorFrom) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	return &Client{
		client:   client,
		store:    store,
		currConv: currConv,
	}, nil
}

func (c *Client) SendMessage(text string, userID int64) error {
	_, err := c.client.Send(tgbotapi.NewMessage(userID, text))
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}
	return nil
}

func (c *Client) setProcUserState(procs []userstateprocessors.UserStateProcessor, state *userstates.UserState) {
	for _, proc := range procs {
		proc.SetUserState(state)
	}
}

func (c *Client) ListenUpdates(msgModel *messages.Model) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	userStateProcessors := []userstateprocessors.UserStateProcessor{
		userstateprocessors.NewCategoryProcessor(),
		userstateprocessors.NewAmountProcessor(c.currConv),
		userstateprocessors.NewDateProcessor(),
		userstateprocessors.NewCurrencyProcessor(c.store.Currency()),
	}

	updates := c.client.GetUpdatesChan(u)
	log.Println("listening for messages")

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			uid := update.Message.From.ID
			text := update.Message.Text
			userState, err := c.store.UserState().GetOne(uid)
			if err == nil && userState.GetStatus() != userstates.ExpectedCommand {
				c.setProcUserState(userStateProcessors, userState)
				for _, proc := range userStateProcessors {
					if userState.GetStatus() == proc.GetProcessStatus() {
						proc.DoProcess(text)
					}
				}
				if userState.Added() {
					err := c.store.Expense().Add(userState.ToExpense())
					if err != nil {
						log.Println("error adding expense:", err)
					}
				}
			}
			if userState == nil {
				userState = userstates.NewUserState(uid)
			}

			newStatus, err := msgModel.IncomingMessage(messages.Message{
				Text:   text,
				UserID: uid,
			}, userState)
			if err != nil {
				log.Println("error processing message:", err)
			}
			userState.SetStatus(newStatus)
			err = c.store.UserState().Save(userState)
			if err != nil {
				log.Println("error saving user state:", err)
			}
		}
	}
}
