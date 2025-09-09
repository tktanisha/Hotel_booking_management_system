package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	instance DB
	once     sync.Once
)

type PostgresDB struct {
	db DB
}

func InitPostgres(dbURL string) (DB, error) {
	once.Do(func() {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			return
		}

		if err := db.Ping(); err != nil {
			log.Fatalf("Failed to ping PostgreSQL: %v", err)
		}

		instance = &PostgresDB{db: db}
		log.Println("PostgreSQL connection established")
	})
	return instance, nil
}

func RunMigrations(db DB, filePath string) error {
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("failed to execute migrations: %v", err)
	}
	fmt.Println("Database migrations executed successfully")
	return nil
}

func (p *PostgresDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return p.db.Query(query, args...)
}

func (p *PostgresDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return p.db.QueryRow(query, args...)
}

func (p *PostgresDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return p.db.Exec(query, args...)
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}
