package config

import (
	"log"
	"os"
	"pack-management/internal/pkg/helpers"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		DBHost     string `env:"DB_HOST,required"`
		DBPort     string `env:"DB_PORT,required"`
		DBName     string `env:"DB_NAME,required"`
		DBUser     string `env:"DB_USER,required"`
		DBPassword string `env:"DB_PASSWORD,required"`
	}
)

func NewConfig() (*Config, error) {
	err := loadEnv()
	if err != nil {
		log.Println("error loading env file")
	}

	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func loadEnv() error {
	goenv := os.Getenv("GO_ENV")
	envFile := ".env"

	if goenv != "" {
		envFile = ".env." + goenv
	}

	rootDir, err := helpers.GetRootDirectory()
	if err != nil {
		return err
	}

	return godotenv.Load(rootDir + "/" + envFile)
}
