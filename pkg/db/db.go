package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Arash-Afshar/gohtmx/pkg/models"
)

func NewDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("openning connection: %v", err)
	}
	if _, err = db.Exec(createSamplesTable); err != nil {
		return nil, fmt.Errorf("creating tables: %v", err)
	}
	if err = CreatePostsTable(context.Background(), db); err != nil {
		return nil, fmt.Errorf("creating tables: %v", err)
	}
	return db, nil
}

func Query[T any, Ptr interface{ *T }](
	ctx context.Context,
	db *sql.DB,
	query string,
	binder func(Ptr) []any,
	args ...any,
) ([]T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []T

	for rows.Next() {
		var result T
		cols := binder(&result)
		err = rows.Scan(cols...)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func QueryScanMany[T any, Ptr models.Scanner[T]](
	ctx context.Context,
	db *sql.DB,
	query string,
	args ...any,
) ([]T, error) {
	return Query(ctx, db, query, func(t Ptr) []any {
		return t.Scan()
	}, args...)
}

func QueryScanOne[T any, Ptr models.Scanner[T]](
	ctx context.Context,
	db *sql.DB,
	query string,
	args ...any,
) (*T, error) {
	rows, err := Query(ctx, db, query, func(t Ptr) []any {
		return t.Scan()
	}, args...)
	if err != nil {
		return nil, fmt.Errorf("querying (%s) with (%v): [%v]", query, args, err)
	}
	if len(rows) != 1 {
		return nil, errors.New("found more than 1 match")
	}
	return &rows[0], nil
}

func Exec(
	ctx context.Context,
	db *sql.DB,
	query string,
	args ...any,
) error {
	_, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("executing %s with (%s): %v", query, args, err)
	}
	return nil
}
