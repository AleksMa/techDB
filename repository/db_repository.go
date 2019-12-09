package repository

import (
	"database/sql"
)

type DBStore struct {
	DB *sql.DB
}

func NewDBStore(db *sql.DB) Repo {
	return &DBStore{
		db,
	}
}
