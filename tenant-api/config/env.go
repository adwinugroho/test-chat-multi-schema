package config

type EnvAppConfig struct {
	AppVersion  string `mapstructure:"app_version"`
	AppName     string `mapstructure:"app_name"`
	Port        string `mapstructure:"app_port"`
	AppURL      string `mapstructure:"app_url"`
	Environment string `mapstructure:"environment"`
	JWTSecret   string `mapstructure:"jwt_secret"`
}
