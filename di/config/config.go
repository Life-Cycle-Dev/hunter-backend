package config

import "github.com/kelseyhightower/envconfig"

type AppConfig struct {
	ServerConfig   ServerConfig
	DatabaseConfig DatabaseConfig
	EmailConfig    EmailConfig
}

type ServerConfig struct {
	Env       string `envconfig:"APP_ENV" default:"local"`
	Port      string `envconfig:"APP_PORT" default:"9000"`
	SentryDns string `envconfig:"APP_SENTRY_DNS" default:"https://sentry.io"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     int    `envconfig:"DB_PORT" default:"3306"`
	User     string `envconfig:"DB_USER" default:"root"`
	Password string `envconfig:"DB_PASSWORD" default:""`
	DBName   string `envconfig:"DB_NAME" default:"hunter"`
}

type EmailConfig struct {
	Host     string `envconfig:"SMTP_HOST" default:"localhost@gmail.com"`
	Port     int    `envconfig:"SMTP_PORT" default:"587"`
	Sender   string `envconfig:"SMTP_SENDER" default:"localhost"`
	Password string `envconfig:"SMTP_PASSWORD" default:"localhost"`
}

func GetConfig() AppConfig {
	var app AppConfig
	envconfig.MustProcess("APP", &app.ServerConfig)
	envconfig.MustProcess("APP", &app.DatabaseConfig)
	envconfig.MustProcess("APP", &app.EmailConfig)
	return app
}
