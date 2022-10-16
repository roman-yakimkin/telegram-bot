package main

import (
	"context"
	"log"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/db/postgres"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo/memrepo"
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
	mode := "postgres"
	var currencyRepo repo.CurrencyRepo
	var currencyRateRepo repo.CurrencyRateRepo
	var expRepo repo.ExpensesRepo
	var limitRepo repo.ExpenseLimitsRepo
	var userStateRepo repo.UserStateRepo
	switch mode {
	case "memory":
		currencyRepo, err = memrepo.NewCurrencyRepo(cfg)
		if err != nil {
			log.Fatal("currency init failed:", err)
		}
		currencyRateRepo = memrepo.NewCurrencyRateRepo(cfg)
		expRepo = memrepo.NewExpenseRepo()
		limitRepo = memrepo.NewExpenseLimitsRepo(cfg)
		userStateRepo = memrepo.NewUserStateRepo()
	case "postgres":
		db := postgres.NewDBConnect(cfg)
		pool, err := db.Connect(ctx)
		if err != nil {
			log.Fatal("Postgres pool connect failed:", err)
		}
		defer db.Disconnect(ctx)
		currencyRepo = pgsqlrepo.NewCurrencyRepo(ctx, pool)
		currencyRateRepo = pgsqlrepo.NewCurrencyRateRepo(ctx, pool, cfg)
		expRepo = pgsqlrepo.NewExpenseRepo(ctx, pool, cfg)
		limitRepo = pgsqlrepo.NewExpenseLimitsRepo(ctx, pool, cfg)
		userStateRepo = pgsqlrepo.NewUserStateRepo(ctx, pool)
	}

	tickerInterval := time.Second * time.Duration(cfg.GetConfig().CurrencyRateLoadInterval)
	daysCount := cfg.GetConfig().CurrencyRateGetDaysCount
	currencyUpdate := tickers.NewCurrencyUpdate(currencyRateRepo, tickerInterval, daysCount)
	currencyUpdate.Run(ctx)

	currencyConvertor := convertors.NewCurrencyConvertor(currencyRateRepo, cfg)
	store := implstore.NewStore(expRepo, userStateRepo, currencyRepo, currencyRateRepo, limitRepo)

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

	if err = tgClient.ListenUpdates(msgModel); err != nil {
		log.Fatal(err)
	}
}
