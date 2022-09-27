package main

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo/memory"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/reports"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	expRepo := memory.New()
	rm := reports.New(expRepo)

	tgClient, err := tg.New(cfg, expRepo, rm)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	msgModel := messages.New(tgClient, rm)

	tgClient.ListenUpdates(msgModel)
}
