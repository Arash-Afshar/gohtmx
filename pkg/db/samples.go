package db

import (
	"context"
	"database/sql"

	"github.com/Arash-Afshar/gohtmx/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

func AddSample(ctx context.Context, db *sql.DB, sample *models.Sample) error {
	return Exec(ctx, db, `INSERT INTO samples(id, name) VALUES (?, ?)`, sample.Id, sample.Name)
}

func DeleteSample(ctx context.Context, db *sql.DB, sample *models.Sample) error {
	return Exec(ctx, db, `DELETE FROM samples WHERE id = ?`, sample.Id)
}

func ListSamples(ctx context.Context, db *sql.DB) ([]models.Sample, error) {
	return QueryScanMany[models.Sample](ctx, db, "SELECT id, name FROM samples")
}

func FindSample(ctx context.Context, db *sql.DB, id string) (*models.Sample, error) {
	return QueryScanOne[models.Sample](ctx, db, "SELECT id, name FROM samples WHERE id = ?", id)
}
