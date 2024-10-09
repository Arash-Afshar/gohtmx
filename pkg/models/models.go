package models

type Scanner[T any] interface {
	*T
	Scan() []any
}
