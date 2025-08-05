package config

import (
	"fmt"
	"log"

	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	"github.com/spf13/viper"
)

var (
	AppConfig EnvAppConfig
)

var configStruct = map[string]interface{}{
	"app-config":        &AppConfig,
	"postgresql-config": &PostgreSQLConfig,
}

func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	for key, value := range configStruct {
		if err := viper.Unmarshal(value); err != nil {
			logger.LogFatal(fmt.Sprintf("Error loading config %s, cause: %+v\n", key, err))
		}
		logger.LogInfo(fmt.Sprintf("Config loaded successfully: %s, value: %+v", key, value))
	}
}
