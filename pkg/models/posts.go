package models

import (
	"github.com/google/uuid"
)

type Scanner[T any] interface {
	*T
	Scan() []any
}

type Post struct {
	Id    string
	Title string
}

func (p *Post) Scan() []any {
	return []any{&p.Id, &p.Title}
}

func NewPost(title string) *Post {
	return &Post{
		Id:    uuid.New().String(),
		Title: title,
	}
}
