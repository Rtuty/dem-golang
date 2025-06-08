package database

import (
	"database/sql"
	"fmt"
	"log"

	"wallpaper-system/internal/config"

	_ "github.com/lib/pq"
)

// DB представляет подключение к базе данных
type DB struct {
	conn *sql.DB
}

// New создает новое подключение к базе данных
func New(cfg *config.DatabaseConfig) (*DB, error) {
	dsn := cfg.GetDSN()

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Проверка подключения
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка проверки подключения к базе данных: %w", err)
	}

	log.Println("Успешное подключение к базе данных")

	return &DB{conn: conn}, nil
}

// Close закрывает подключение к базе данных
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// GetConnection возвращает подключение к базе данных
func (db *DB) GetConnection() *sql.DB {
	return db.conn
}

// Ping проверяет подключение к базе данных
func (db *DB) Ping() error {
	return db.conn.Ping()
}
