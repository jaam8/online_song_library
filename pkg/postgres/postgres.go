package postgres

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host             string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" env-default:"localhost"`
	Port             string `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" env-default:"5432"`
	Username         string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" env-default:"root"`
	Password         string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD" env-default:"1234"`
	Database         string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" env-default:"postgres"`
	PathToMigrations string `yaml:"PATH_TO_MIGRATIONS" env:"PATH_TO_MIGRATIONS" env-default:"file:///app/db/migrations"`
}

// New создает подключение к postgres
func New(config Config) (*gorm.DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	conn, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = runMigrations(connString, config.PathToMigrations)
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return conn, nil
}

// runMigrations выполняет миграции
func runMigrations(databaseURL, migrationsPath string) error {
	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}
