package db

import (
	"database/sql"
	"fmt"
)

func NewDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("openning connection: %v", err)
	}
	if _, err = db.Exec(createSamplesTable); err != nil {
		return nil, fmt.Errorf("creating tables: %v", err)
	}
	if _, err = db.Exec(createPostTable); err != nil {
		return nil, fmt.Errorf("creating tables: %v", err)
	}
	return db, nil
}
