package config

import (
	"github.com/gookit/slog"
	"github.com/spf13/viper"
	"time"
)

var Config config

type config struct {
	Server   Server
	Postgres Postgres
	JWT      JWT
}

type Server struct {
	Host string
	Port string
}

type Postgres struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

type JWT struct {
	JWTSecret   string
	TokenExpiry time.Duration
}

func init() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		slog.Errorf("Ошибка при чтении конфигурации: %s", err)
	}

	tokenExpiry, err := time.ParseDuration(viper.GetString("TOKEN_EXPIRY"))
	if err != nil {
		slog.Errorf("Ошибка при парсинге длительности токена: %s", err)
		tokenExpiry = time.Hour * 24
	}

	Config = config{
		Server: Server{
			Host: viper.GetString("SRV_HOST"),
			Port: viper.GetString("SRV_PORT"),
		},
		Postgres: Postgres{
			Username: viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetString("POSTGRES_PORT"),
			DBName:   viper.GetString("POSTGRES_DB"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		JWT: JWT{
			JWTSecret:   viper.GetString("SECRET_KEY"),
			TokenExpiry: tokenExpiry,
		},
	}

}
