package models

type Sample struct {
	Id   int64
	Name string
}

func NewSample(name string) *Sample {
	return &Sample{
		Id:   0,
		Name: name,
	}
}
