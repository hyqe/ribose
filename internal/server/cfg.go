package server

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port          string `default:"80"`
	Env           string `default:"DEV"` // DEV, PROD
	PostgresURL   string `split_words:"true" required:"true"`
	MigrationsURL string `split_words:"true" required:"true"`
}

func loadConfig() (c Config, err error) {
	return c, envconfig.Process("", &c)
}

func (c *Config) Addr() string {
	return fmt.Sprintf(":%v", c.Port)
}
