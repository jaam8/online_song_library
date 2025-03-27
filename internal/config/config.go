package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jaam8/online_song_library/pkg/postgres"
	"github.com/joho/godotenv"
)

type Config struct {
	RestPort   string          `yaml:"REST_PORT" env:"REST_PORT" env-default:"8080"`
	SwaggerUrl string          `yaml:"SWAGGER_URL" env:"SWAGGER_URL"`
	LogLevel   string          `yaml:"LOG_LEVEL" env:"LOG_LEVEL" env-default:"debug"`
	Postgres   postgres.Config `yaml:"POSTGRES" env:"POSTGRES"`
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	var config Config
	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
