package main

import (
	"context"
	"flag"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/db/postgres"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/importers"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/localmetrics"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/localtracing"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/loggers"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo/pgsqlrepo"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/store/implstore"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/tickers"
	"go.uber.org/zap"
)

var develMode = flag.Bool("devel", false, "development mode")

func main() {
	logger := loggers.InitLogger(*develMode)
	localtracing.InitTracing(logger)

	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		logger.Fatal("config init failed:", zap.Error(err))
	}

	db := postgres.NewDBConnect(cfg)
	pool, err := db.Connect(ctx)
	if err != nil {
		logger.Fatal("Postgres pool connect failed:", zap.Error(err))
	}
	defer db.Disconnect(ctx)

	currencyRateImporter := importers.NewCbrRateImporter(pool, cfg)
	currencyRepo := pgsqlrepo.NewCurrencyRepo(pool)
	currencyRateRepo := pgsqlrepo.NewCurrencyRateRepo(pool, currencyRateImporter)

	expRepo := pgsqlrepo.NewExpenseRepo(pool, cfg, logger)
	limitRepo := pgsqlrepo.NewExpenseLimitsRepo(pool, cfg)
	userStateRepo := pgsqlrepo.NewUserStateRepo(pool)

	if err := userStateRepo.ClearStatus(ctx); err != nil {
		logger.Fatal("error clear user status: ", zap.Error(err))
	}

	tickerInterval := time.Second * time.Duration(cfg.GetConfig().CurrencyRateLoadInterval)
	daysCount := cfg.GetConfig().CurrencyRateGetDaysCount
	currencyUpdate := tickers.NewCurrencyUpdate(currencyRateRepo, tickerInterval, daysCount, logger)
	currencyUpdate.Run(ctx)

	currencyConvertor := convertors.NewCurrencyConvertor(currencyRateRepo, cfg)
	store := implstore.NewStore(expRepo, userStateRepo, currencyRepo, currencyRateRepo, limitRepo, currencyConvertor)

	currencyOutput := output.NewCurrencyListOutput(currencyRepo)
	currencyAmountOutput := output.NewCurrencyAmount(currencyRepo)
	reportManager := output.NewReportManager(store, currencyConvertor, currencyAmountOutput, logger)
	limitListOutput := output.NewLimitListOutput(limitRepo, currencyAmountOutput)
	outputSet := output.NewOutput(currencyOutput, reportManager, limitListOutput)

	tgClient, err := tg.New(cfg, store, currencyConvertor, logger)
	if err != nil {
		logger.Fatal("tg client init failed:", zap.Error(err))
	}

	msgModel := messages.New(tgClient, outputSet, logger)

	go localmetrics.HandleMetrics(cfg.GetConfig(), logger)

	if err = tgClient.ListenUpdates(ctx, msgModel); err != nil {
		logger.Fatal("tg client error listening updates", zap.Error(err))
	}
}
