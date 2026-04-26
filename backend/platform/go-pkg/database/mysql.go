package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLConfig konfigurasi koneksi MySQL
type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	// Pool settings
	MaxOpenConns    int           // Maksimal koneksi terbuka
	MaxIdleConns    int           // Maksimal koneksi idle
	ConnMaxLifetime time.Duration // Umur maksimal koneksi
}

// NewMySQL buat koneksi MySQL dengan pool
func NewMySQL(cfg MySQLConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql: %w", err)
	}

	// Konfigurasi pool
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	} else {
		db.SetMaxOpenConns(25) // Default
	}

	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	} else {
		db.SetMaxIdleConns(5) // Default
	}

	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	} else {
		db.SetConnMaxLifetime(5 * time.Minute) // Default
	}

	// Test koneksi
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	return db, nil
}

// HealthCheck cek koneksi database masih hidup
func HealthCheck(db *sql.DB) error {
	return db.Ping()
}