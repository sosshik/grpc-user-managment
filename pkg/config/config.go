package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port        string `env:"PORT"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
	DbUrl       string `env:"DATABASE_URL"`
	ReconnTime  int    `env:"RECONN_TIME" envDefault:"5"`
	ConnCheck   bool   `env:"CONN_CHECK" envDefault:"true"`
	ReconnTries int    `env:"RECONN_TRIES" envDefault:"5"`
}

var once sync.Once

var configInstance *Config

func GetConfig() *Config {
	if configInstance == nil {
		once.Do(func() {
			fmt.Println("Creating config instance now.")

			var cfg Config

			if err := env.Parse(&cfg); err != nil {
				log.Fatal(err)
			}

			configInstance = &cfg
		})

	}

	return configInstance
}
