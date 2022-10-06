package main

import (
	"context"
	"log"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo/memrepo"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/store/implstore"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/tickers"
)

func main() {
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	expRepo := memrepo.NewExpenseRepo()
	userStateRepo := memrepo.NewUserStateRepo()
	currencyRepo := memrepo.NewCurrencyRepo(cfg)

	tickerInterval := time.Second * time.Duration(cfg.GetConfig().CurrencyRateLoadInterval)
	currencyUpdate := tickers.NewCurrencyUpdate(currencyRepo, tickerInterval)
	currencyUpdate.Run(ctx)

	err = currencyRepo.LoadAll()
	if err != nil {
		log.Fatal("cannot load currency rates", err)
	}

	store := implstore.NewStore(expRepo, userStateRepo, currencyRepo)

	currencyOutput := output.NewCurrencyOutput(currencyRepo)
	currencyConvertor := convertors.NewCurrency(currencyRepo, cfg)
	currencyAmountOutput := output.NewCurrencyAmount(currencyRepo)

	reportManager := output.NewReportManager(store, currencyConvertor, currencyAmountOutput)
	outputSet := output.NewOutput(currencyOutput, reportManager)

	tgClient, err := tg.New(cfg, store, currencyConvertor)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	msgModel := messages.New(tgClient, outputSet)

	tgClient.ListenUpdates(msgModel)
}
