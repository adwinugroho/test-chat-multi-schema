package config

import (
	"context"
	"time"

	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type EnvPostgreSQLConfig struct {
	Database DatabaseConfig `mapstructure:"database"`
}

type DatabaseConfig struct {
	URL string `mapstructure:"url"`
}

var (
	PostgreSQLConfig EnvPostgreSQLConfig
)

type (
	PostgresDB struct {
		DB *pgxpool.Pool
	}
)

func InitConnectDB(ctx context.Context, conn string) (*PostgresDB, error) {
	logger.LogWithFields(logrus.Fields{
		"info": "Connecting to database",
		"url":  conn,
	}, "info connecting to database")

	poolConfig, err := pgxpool.ParseConfig(conn)
	if err != nil {
		logger.LogFatal("Failed to parse config:" + err.Error())
	}

	poolConfig.HealthCheckPeriod = time.Minute * 5
	poolConfig.MaxConns = 100
	poolConfig.MaxConnIdleTime = time.Minute * 1
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeDescribeExec
	poolConfig.ConnConfig.RuntimeParams = map[string]string{}
	poolConfig.ConnConfig.RuntimeParams["application_name"] = "tenant-api-service"
	dbPool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.LogFatal("Failed to create pool:" + err.Error())
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		logger.LogFatal("Failed to ping:" + err.Error())
	}

	return &PostgresDB{DB: dbPool}, nil
}
