package config

import (
	"fmt"

	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	"github.com/spf13/viper"
)

var (
	AppConfig EnvAppConfig
)

var configStruct = map[string]interface{}{
	"app-config":        &AppConfig,
	"postgresql-config": &PostgreSQLConfig,
	"rabbitmq-config":   &RabbitMQConfig,
}

func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.LogFatal("Error while read in config: " + err.Error())
	}

	for key, value := range configStruct {
		if err := viper.Unmarshal(value); err != nil {
			logger.LogFatal(fmt.Sprintf("Error loading config %s, cause: %+v\n", key, err))
		}
	}
}
