package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server
	Postgres
	Session
}

type Server struct {
	SERVER_PORT          int           `envconfig:"SERVER_PORT"`
	SERVER_IDLE_TIMEOUT  time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"5m"`
	SERVER_READ_TIMEOUT  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"10s"`
	SERVER_WRITE_TIMEOUT time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"10s"`
}

type Postgres struct {
	POSTGRES_USERNAME string `envconfig:"POSTGRES_USERNAME"`
	POSTGRES_PASSWORD string `envconfig:"POSTGRES_PASSWORD"`
	POSTGRES_HOST     string `envconfig:"POSTGRES_HOST"`
	POSTGRES_PORT     int    `envconfig:"POSTGRES_PORT"`
	POSTGRES_NAME     string `envconfig:"POSTGRES_NAME"`

	POSTGRES_SSLMODE      string `envconfig:"POSTGRES_SSLMODE"`
	POSTGRES_MAX_CONNECTS int    `envconfig:"POSTGRES_MAX_CONNECTS" default:"10"`
}

type Session struct {
	SESSION_MAX_USER_SESSIONS int `envconfig:"SESSION_MAX_USER_SESSIONS" default:"5"`

	SESSION_EXP_SHORT         time.Duration `envconfig:"SESSION_EXP_SHORT" default:"15m"`
	SESSION_REFRESH_EXP_SHORT time.Duration `envconfig:"SESSION_REFRESH_EXP_SHORT" default:"24h"`

	SESSION_EXP         time.Duration `envconfig:"SESSION_EXP" default:"1h"`
	SESSION_REFRESH_EXP time.Duration `envconfig:"SESSION_REFRESH_EXP" default:"168h"`

	SESSION_EXP_LONG        time.Duration `envconfig:"SESSION_EXP_LONG" default:"24h"`
	SESSION_REFRESH_EXPLONG time.Duration `envconfig:"SESSION_REFRESH_EXPLONG" default:"8766h"`
}

func Load(filenames ...string) (Config, error) {
	if err := godotenv.Load(filenames...); err != nil {
		return Config{}, err
	}

	return Get()
}

func Get() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)

	return cfg, err
}

func PgUrl(cfg Config) string {
	url := fmt.Sprintf(`postgres://%s:%s@%s:%d/%s`,
		cfg.POSTGRES_USERNAME,
		cfg.POSTGRES_PASSWORD,
		cfg.POSTGRES_HOST,
		cfg.POSTGRES_PORT,
		cfg.POSTGRES_NAME,
	)
	if cfg.POSTGRES_SSLMODE == "disable" {
		url = fmt.Sprintf("%s?sslmode=disable",
			url,
		)
	}

	return url
}
