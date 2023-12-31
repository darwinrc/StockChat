package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DB interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	user, password, database, host := os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_HOST")

	ds := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, password, host, database)

	db, err := sql.Open("postgres", ds)
	if err != nil {
		return nil, err
	}

	//return &database{db: db}, nil
	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
