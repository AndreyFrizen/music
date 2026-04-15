package migrator

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/pressly/goose/v3"
)

// DatabaseConfig содержит конфигурацию для одной базы данных
type DatabaseConfig struct {
	Name          string
	DSN           string
	MigrationsDir string
}

type Migrator struct {
	logger  *slog.Logger
	configs []DatabaseConfig
}

func NewMigrator(logger *slog.Logger, configs []DatabaseConfig) *Migrator {
	return &Migrator{
		logger:  logger,
		configs: configs,
	}
}

// Run выполняет миграции для всех баз данных
func (m *Migrator) Run() error {
	for _, cfg := range m.configs {
		if err := m.migrateDatabase(cfg); err != nil {
			return fmt.Errorf("migration failed for %s: %w", cfg.Name, err)
		}
	}
	return nil
}

func (m *Migrator) migrateDatabase(cfg DatabaseConfig) error {
	m.logger.Info("Starting migrations", "database", cfg.Name, "dir", cfg.MigrationsDir)

	// Открываем соединение с БД
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	// Устанавливаем timeout для миграций (опционально)
	goose.SetBaseFS(nil) // если используете embed
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	// Запускаем миграции
	if err := goose.Up(db, cfg.MigrationsDir); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}

	m.logger.Info("Migrations completed", "database", cfg.Name)
	return nil
}
