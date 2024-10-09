package models

import (
	"github.com/google/uuid"
)

type Post struct {
	Id    string
	Title string
}

func NewPost(title string) *Post {
	return &Post{
		Id:    uuid.New().String(),
		Title: title,
	}
}
