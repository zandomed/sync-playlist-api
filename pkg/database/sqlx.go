package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/zandomed/sync-playlist-api/internal/config"
)

// DB encapsula la conexión SQLX
type DB struct {
	*sqlx.DB
}

// Connect establece conexión con PostgreSQL usando SQLX
func Connect(cfg *config.DatabaseConfig) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configurar connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verificar conexión
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{DB: db}, nil
}

// HealthCheck verifica la conexión a la base de datos
func (db *DB) HealthCheck() error {
	return db.Ping()
}

// Transaction ejecuta una función dentro de una transacción
func (db *DB) Transaction(fn func(*sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

// Named wrapper para queries con parámetros nombrados
func (db *DB) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return db.DB.NamedQuery(query, arg)
}

// NamedExec wrapper para ejecución con parámetros nombrados
func (db *DB) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return db.DB.NamedExec(query, arg)
}
