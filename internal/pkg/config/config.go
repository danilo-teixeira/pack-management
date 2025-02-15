package config

import "github.com/caarlos0/env/v11"

type (
	config struct {
		DBHost     string `env:"DB_HOST,required"`
		DBPort     string `env:"DB_PORT,required"`
		DBName     string `env:"DB_NAME,required"`
		DBUser     string `env:"DB_USER,required"`
		DBPassword string `env:"DB_PASSWORD,required"`
	}
)

func NewConfig() (*config, error) {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
