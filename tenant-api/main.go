package main

import (
	"context"
	"log"
	"time"

	"github.com/adwinugroho/test-chat-multi-schema/config"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	"github.com/sirupsen/logrus"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logger.InitLogger()

	config.LoadConfig()
}

func main() {
	logger.LogInfo("Starting application...")
	logger.LogInfo("Application started on port:" + config.AppConfig.Port)
	logger.LogInfo("Application URL:" + config.AppConfig.AppURL)

	parentCtx := context.Background()
	ctx, cancel := context.WithTimeout(parentCtx, 60*time.Second)
	defer cancel()

	dbHandler, err := config.InitConnectDB(ctx, config.PostgreSQLConfig.Database.URL)
	if err != nil {
		logger.LogFatal("Failed to connect to database:" + err.Error())
	}

	logger.LogWithFields(logrus.Fields{
		"db-stat": dbHandler.DB.Stat(),
	}, "checking stat db")
}
