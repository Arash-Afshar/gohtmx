package models

import "github.com/google/uuid"

type Sample struct {
	Id   string
	Name string
}

func NewSample(name string) *Sample {
	return &Sample{
		Id:   uuid.New().String(),
		Name: name,
	}
}

func (s *Sample) Scan() []any {
	return []any{&s.Id, &s.Name}
}
