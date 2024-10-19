package models

import (
	"github.com/google/uuid"
)

type Post struct {
	Id      string
	Title   string
	Content string
}

func (p *Post) Scan() []any {
	return []any{&p.Id, &p.Title, &p.Content}
}

func NewPost(title string, content string) *Post {
	return &Post{
		Id:      uuid.New().String(),
		Title:   title,
		Content: content,
	}
}
