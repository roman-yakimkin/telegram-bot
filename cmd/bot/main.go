package main

import (
	"context"
	"log"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/db/postgres"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/importers"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo/pgsqlrepo"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/store/implstore"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/tickers"
)

func main() {
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	db := postgres.NewDBConnect(cfg)
	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatal("Postgres pool connect failed:", err)
	}
	defer db.Disconnect(ctx)

	currencyRateImporter := importers.NewCbrRateImporter(pool, cfg)
	currencyRepo := pgsqlrepo.NewCurrencyRepo(pool)
	currencyRateRepo := pgsqlrepo.NewCurrencyRateRepo(pool, currencyRateImporter)

	expRepo := pgsqlrepo.NewExpenseRepo(pool, cfg)
	limitRepo := pgsqlrepo.NewExpenseLimitsRepo(pool, cfg)
	userStateRepo := pgsqlrepo.NewUserStateRepo(pool)

	if err := userStateRepo.ClearStatus(ctx); err != nil {
		log.Fatal("error clear user status: ", err)
	}

	tickerInterval := time.Second * time.Duration(cfg.GetConfig().CurrencyRateLoadInterval)
	daysCount := cfg.GetConfig().CurrencyRateGetDaysCount
	currencyUpdate := tickers.NewCurrencyUpdate(currencyRateRepo, tickerInterval, daysCount)
	currencyUpdate.Run(ctx)

	currencyConvertor := convertors.NewCurrencyConvertor(currencyRateRepo, cfg)
	store := implstore.NewStore(expRepo, userStateRepo, currencyRepo, currencyRateRepo, limitRepo, currencyConvertor)

	currencyOutput := output.NewCurrencyListOutput(currencyRepo)
	currencyAmountOutput := output.NewCurrencyAmount(currencyRepo)
	reportManager := output.NewReportManager(store, currencyConvertor, currencyAmountOutput)
	limitListOutput := output.NewLimitListOutput(limitRepo, currencyAmountOutput)
	outputSet := output.NewOutput(currencyOutput, reportManager, limitListOutput)

	tgClient, err := tg.New(cfg, store, currencyConvertor)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	msgModel := messages.New(tgClient, outputSet)

	if err = tgClient.ListenUpdates(ctx, msgModel); err != nil {
		log.Fatal(err)
	}
}
