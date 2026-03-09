package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	config "user-service/config/grpc_server"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// DB – обёртка над пулом PostgreSQL и клиентом Redis
type DB struct {
	pg    *pgxpool.Pool
	redis *redis.Client
	log   *slog.Logger
	cfg   *config.Config
}

func NewDB(log *slog.Logger, cfg *config.Config) (*DB, error) {
	const op = "database.NewDB"

	if cfg == nil {
		return nil, fmt.Errorf("%s: config is nil", op)
	}

	// Настройка пула PostgreSQL
	pgConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse postgres DSN: %w", op, err)
	}

	// Дополнительные параметры пула (можно вынести в конфиг)
	if cfg.MaxConns > 0 {
		pgConfig.MaxConns = cfg.MaxConns
	}
	if cfg.MinConns > 0 {
		pgConfig.MinConns = cfg.MinConns
	}
	if cfg.MaxConnIdle > 0 {
		pgConfig.MaxConnIdleTime = cfg.MaxConnIdle
	}
	if cfg.ConnTimeout > 0 {
		pgConfig.ConnConfig.ConnectTimeout = cfg.ConnTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgPool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create postgres pool: %w", op, err)
	}

	// Проверяем соединение
	if err := pgPool.Ping(ctx); err != nil {
		pgPool.Close()
		return nil, fmt.Errorf("%s: postgres ping failed: %w", op, err)
	}
	log.Info("Connected to PostgreSQL")

	// Инициализация Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		pgPool.Close()
		return nil, fmt.Errorf("%s: redis ping failed: %w", op, err)
	}
	log.Info("Connected to Redis", "addr", cfg.RedisHost)

	return &DB{
		pg:    pgPool,
		redis: redisClient,
		log:   log,
		cfg:   cfg,
	}, nil
}

func (d *DB) Close() error {
	var errs []error

	if d.pg != nil {
		d.pg.Close()
		d.log.Info("PostgreSQL connection closed")
	}

	if d.redis != nil {
		if err := d.redis.Close(); err != nil {
			errs = append(errs, fmt.Errorf("redis close error: %w", err))
		} else {
			d.log.Info("Redis connection closed")
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing DB: %v", errs)
	}
	return nil
}

func (d *DB) Ping(ctx context.Context) error {
	if err := d.pg.Ping(ctx); err != nil {
		return fmt.Errorf("postgres ping failed: %w", err)
	}
	if err := d.redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}

func (d *DB) GetPG() *pgxpool.Pool {
	return d.pg
}

func (d *DB) GetRedis() *redis.Client {
	return d.redis
}

func (db *DB) QueryRowContext(ctx context.Context, query string, args ...any) pgx.Row {
	return db.pg.QueryRow(ctx, query, args...)
}

// ExecContext выполняет запрос без возврата строк.
func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return db.pg.Exec(ctx, query, args...)
}

// QueryContext выполняет запрос, возвращая несколько строк.
func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return db.pg.Query(ctx, query, args...)
}
