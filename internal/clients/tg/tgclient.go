package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/interfaces"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/expenses"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/reports"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/vars"
	"log"
	"strconv"
	"time"
)

type TokenGetter interface {
	Token() string
}

type UserInputData struct {
	status        int
	category      string
	amount        int
	date          time.Time
	addedCategory bool
	addedAmount   bool
	addedDate     bool
}

func (ui *UserInputData) Added() bool {
	return ui.addedCategory && ui.addedAmount && ui.addedDate
}

type Client struct {
	client  *tgbotapi.BotAPI
	expRepo interfaces.ExpensesRepo
	rm      *reports.ReportManager
	inpData map[int64]UserInputData
}

func New(tokenGetter TokenGetter, er interfaces.ExpensesRepo, rm *reports.ReportManager) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	return &Client{
		client:  client,
		expRepo: er,
		rm:      rm,
		inpData: make(map[int64]UserInputData),
	}, nil
}

func (c *Client) SendMessage(text string, userID int64) error {
	_, err := c.client.Send(tgbotapi.NewMessage(userID, text))
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}
	return nil
}

func (c *Client) ListenUpdates(msgModel *messages.Model) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	log.Println("listening for messages")

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			uid := update.Message.From.ID
			text := update.Message.Text
			currentStatus := vars.ExpectedCommand
			inp, ok := c.inpData[uid]
			if ok {
				switch {
				case inp.status == vars.ExpectedCategory:
					inp.category = text
					inp.addedCategory = true
				case inp.status == vars.ExpectedAmount:
					var err error
					inp.amount, err = strconv.Atoi(text)
					if err != nil {
						inp.status = vars.IncorrectAmount
						break
					}
					inp.addedAmount = true
				case inp.status == vars.ExpectedDate:
					if text == "*" {
						inp.date = time.Now()
						inp.addedDate = true
						break
					}
					var err error
					inp.date, err = time.Parse("2006-01-02", text)
					if err != nil {
						inp.status = vars.IncorrectDate
						break
					}
					inp.addedDate = true
				}
				if inp.Added() {
					err := c.expRepo.Add(&expenses.Expense{
						UserID:   uid,
						Category: inp.category,
						Amount:   inp.amount,
						Date:     inp.date,
					})
					if err != nil {
						log.Println("error adding expense:", err)
					}
				}
				currentStatus = inp.status
			}

			newStatus, err := msgModel.IncomingMessage(messages.Message{
				Text:   text,
				UserID: uid,
			}, currentStatus)
			if err != nil {
				log.Println("error processing message:", err)
				delete(c.inpData, uid)
			}
			if newStatus == vars.ExpectedCommand {
				delete(c.inpData, uid)
			} else {
				inp.status = newStatus
				c.inpData[uid] = inp
			}
		}
	}
}
