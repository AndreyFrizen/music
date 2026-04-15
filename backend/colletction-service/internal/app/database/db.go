package database

import (
	config "collection-service/config/grpc_server"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type DB struct {
	pg    *pgxpool.Pool
	redis *redis.Client
	log   *slog.Logger
	cfg   *config.Config
}

func NewDB(log *slog.Logger, cfg *config.Config) (*DB, error) {
	op := "database.NewDB"

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=10",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	pgConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse postgres DSN: %w", op, err)
	}

	pgConfig.MaxConnLifetime = 5 * time.Minute
	pgConfig.MaxConnIdleTime = 1 * time.Minute
	pgConfig.HealthCheckPeriod = 30 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgPool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create connection pool: %w", op, err)
	}

	if err := pgPool.Ping(ctx); err != nil {
		pgPool.Close()
		return nil, fmt.Errorf("%s: postgres ping failed: %w", op, err)
	}

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

func (db *DB) QueryRowContext(ctx context.Context, query string, args ...any) pgx.Row {
	return db.pg.QueryRow(ctx, query, args...)
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return db.pg.Exec(ctx, query, args...)
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return db.pg.Query(ctx, query, args...)
}

func (db *DB) GetRedis(ctx context.Context, key string) *redis.StringCmd {
	return db.redis.Get(ctx, key)
}

func (db *DB) SetRedis(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	return db.redis.Set(ctx, key, value, expiration)
}

func (db *DB) DelRedis(ctx context.Context, keys ...string) *redis.IntCmd {
	return db.redis.Del(ctx, keys...)
}

func (db *DB) Log() *slog.Logger {
	return db.log
}
