package config

import "github.com/caarlos0/env/v6"

type Config struct {
	Env        string `env:"TODO_ENV" envDefault:"dev"`
	Port       int    `env:"PORT" envDefault:"80"`
	DBHost     string `env:"DB_HOST" envDefault:"db"`
	DBPort     int    `env:"DB_PORT" envDefault:"13306"`
	DBUser     string `env:"DB_USER" envDefault:"gotodo"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"gotodo"`
	DBName     string `env:"DB_NAME" envDefault:"gotodo"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
