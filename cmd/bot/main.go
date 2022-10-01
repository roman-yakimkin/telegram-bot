package main

import (
	"log"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo/memory"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/reports"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	expRepo := memory.NewExpenseRepo()
	userStateRepo := memory.NewUserStateRepo()
	reportManager := reports.New(expRepo)

	tgClient, err := tg.New(cfg, expRepo, userStateRepo, reportManager)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	msgModel := messages.New(tgClient, reportManager)

	tgClient.ListenUpdates(msgModel)
}
