package copmopites

import "database/sql"

type Store struct {
	db *sql.DB
}

func MSSQLComposite(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}
