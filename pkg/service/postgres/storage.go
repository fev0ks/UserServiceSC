package postgres

import (
	"database/sql"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db}
}
