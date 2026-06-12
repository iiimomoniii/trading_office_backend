package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"trading-office/trading_office_backend/config"
)

func dsn(cfg config.DBConfig) string {
	fmt.Printf("[db] dsn => host=%s port=%d dbname=%s user=%s sslmode=%s\n",
		cfg.Host, cfg.Port, cfg.DBName, cfg.User, cfg.SSLMode)

	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.DBName, cfg.User, cfg.Password, cfg.SSLMode,
	)
}

// NewDatabase — init PostgreSQL connection pool
func NewDatabase(cfg config.DBConfig) (*sql.DB, error) {
	conn, err := sql.Open("postgres", dsn(cfg))
	if err != nil {
		return nil, fmt.Errorf("[db] failed to open connection: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("[db] failed to ping database: %w", err)
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(10)

	return conn, nil
}

// RunMigrations — create schema ก่อน แล้วค่อย run golang-migrate
func RunMigrations(cfg config.DBConfig) error {
	conn, err := sql.Open("postgres", dsn(cfg))
	if err != nil {
		return fmt.Errorf("[db] migration open error: %w", err)
	}
	defer conn.Close()

	// สร้าง schema ก่อนเสมอ — migrate จะใช้ schema นี้ track version
	if _, err := conn.Exec(`CREATE SCHEMA IF NOT EXISTS trading_office`); err != nil {
		return fmt.Errorf("[db] create schema error: %w", err)
	}
	fmt.Println("[db] schema trading_office ready")

	driver, err := postgres.WithInstance(conn, &postgres.Config{
		SchemaName: "trading_office",
	})
	if err != nil {
		return fmt.Errorf("[db] migration driver error: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("[db] migration init error: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("[db] migration failed: %w", err)
	}

	fmt.Println("[db] migrations applied successfully")
	return nil
}
