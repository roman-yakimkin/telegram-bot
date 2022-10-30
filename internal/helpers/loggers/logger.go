package loggers

import (
	"log"

	"go.uber.org/zap"
)

func InitLogger(develMode bool) *zap.Logger {
	var logger *zap.Logger
	var err error
	if develMode {
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatal("init dev logger error:", err)
		}
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			log.Fatal("init dev logger error:", err)
		}
	}
	return logger
}
