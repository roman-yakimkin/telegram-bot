package tg

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/reports"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/services/userstateprocessors"
)

type TokenGetter interface {
	Token() string
}

type Client struct {
	client        *tgbotapi.BotAPI
	expRepo       repo.ExpensesRepo
	userStateRepo repo.UserStateRepo
	rm            *reports.ReportManager
}

func New(tokenGetter TokenGetter, er repo.ExpensesRepo, usr repo.UserStateRepo, rm *reports.ReportManager) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	return &Client{
		client:        client,
		expRepo:       er,
		userStateRepo: usr,
		rm:            rm,
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
		userstateprocessors.NewAmountProcessor(),
		userstateprocessors.NewDateProcessor(),
	}

	updates := c.client.GetUpdatesChan(u)
	log.Println("listening for messages")

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			uid := update.Message.From.ID
			text := update.Message.Text
			currentStatus := userstates.ExpectedCommand
			userState, err := c.userStateRepo.GetOne(uid)
			if err == nil {
				c.setProcUserState(userStateProcessors, userState)
				for _, proc := range userStateProcessors {
					if userState.GetStatus() == proc.GetProcessStatus() {
						proc.DoProcess(text)
					}
				}
				if userState.Added() {
					err := c.expRepo.Add(userState.ToExpense())
					if err != nil {
						log.Println("error adding expense:", err)
					}
				}
				currentStatus = userState.GetStatus()
			}

			newStatus, err := msgModel.IncomingMessage(messages.Message{
				Text:   text,
				UserID: uid,
			}, currentStatus)
			if err != nil {
				log.Println("error processing message:", err)
				err = c.userStateRepo.Delete(uid)
				if err != nil {
					log.Println("error deleting user state:", err)
				}
			}
			if newStatus == userstates.ExpectedCommand {
				err = c.userStateRepo.Delete(uid)
				if err != nil {
					log.Println("error deleting user state:", err)
				}
			} else {
				if userState == nil {
					userState = &userstates.UserState{UserID: uid}
				}
				userState.SetStatus(newStatus)
				err = c.userStateRepo.Save(userState)
				if err != nil {
					log.Println("error saving user state:", err)
				}
			}
		}
	}
}
