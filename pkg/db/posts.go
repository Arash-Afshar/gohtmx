package db

import (
	"context"
	"database/sql"

	"github.com/Arash-Afshar/gohtmx/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

func AddPost(ctx context.Context, db *sql.DB, post *models.Post) error {
	return Exec(ctx, db, `INSERT INTO posts(id, title) VALUES (?, ?)`, post.Id, post.Title)
}

func DeletePost(ctx context.Context, db *sql.DB, post *models.Post) error {
	return Exec(ctx, db, `DELETE FROM posts WHERE id = ?`, post.Id)
}

func ListPosts(ctx context.Context, db *sql.DB) ([]models.Post, error) {
	return QueryScanMany[models.Post](ctx, db, "SELECT id, title FROM posts")
}

func FindPost(ctx context.Context, db *sql.DB, id string) (*models.Post, error) {
	return QueryScanOne[models.Post](ctx, db, "SELECT id, title FROM posts WHERE id = ?", id)
}
