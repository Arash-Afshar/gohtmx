package db

import (
	"database/sql"
	"fmt"

	"github.com/Arash-Afshar/gohtmx/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

const createPostTable string = `
CREATE TABLE IF NOT EXISTS "posts" (
	"id"   TEXT NOT NULL UNIQUE,
	"title" TEXT
);
`

// func Query[T any, Ptr interface{ *T }](
// 	ctx context.Context,
// 	db *sql.DB,
// 	query string,
// 	binder func(Ptr) []any, args ...any,
// ) ([]T, error) {
// 	rows, err := db.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	var results []T
//
// 	for rows.Next() {
// 		var result T
// 		cols := binder(&result)
// 		err = rows.Scan(cols...)
// 		if err != nil {
// 			return nil, err
// 		}
// 		results = append(results, result)
// 	}
// 	return results, nil
// }

func AddPost(db *sql.DB, post *models.Post) error {
	sqlString := `INSERT INTO posts(id, title) VALUES (?, ?)`
	statement, err := db.Prepare(sqlString)
	if err != nil {
		return fmt.Errorf("preparing %s: %v", sqlString, err)
	}
	_, err = statement.Exec(post.Id, post.Title)
	if err != nil {
		return fmt.Errorf("executing %s with (%s): %v", sqlString, post.Title, err)
	}
	return nil
}

func FindPost(db *sql.DB, id string) (*models.Post, error) {
	sqlString := `SELECT title FROM posts WHERE id = ?`
	statement, err := db.Prepare(sqlString)
	if err != nil {
		return nil, fmt.Errorf("preparing %s: %v", sqlString, err)
	}
	row, err := statement.Query(id)
	defer row.Close()
	if row.Next() {
		var title string
		if err = row.Scan(&title); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		if row.Next() {
			return nil, fmt.Errorf("more than one entry found")
		}
		return &models.Post{Id: id, Title: title}, nil
	} else {
		return nil, fmt.Errorf("not found")
	}
}

func DeletePost(db *sql.DB, post *models.Post) error {
	sqlString := `DELETE FROM posts WHERE id = ?`
	statement, err := db.Prepare(sqlString)
	if err != nil {
		return fmt.Errorf("preparing %s: %v", sqlString, err)
	}
	_, err = statement.Exec(post.Id)
	if err != nil {
		return fmt.Errorf("executing %s with (%s): %v", sqlString, post.Id, err)
	}
	return nil
}

func ListPosts(db *sql.DB) ([]*models.Post, error) {
	sqlString := `SELECT id, title FROM posts`
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, fmt.Errorf("preparing %s: %v", sqlString, err)
	}
	defer rows.Close()
	posts := make([]*models.Post, 0)
	for rows.Next() {
		var id string
		var title string
		err = rows.Scan(&id, &title)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		posts = append(posts, &models.Post{
			Id:    id,
			Title: title,
		})
	}
	return posts, nil
}
