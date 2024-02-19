package storage

import (
	"time"
)

type Tasks []Task

type Task struct {
	ID         int
	Opened     time.Time
	Closed     time.Time
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

type SearchOptions struct {
	FilterOptions
	PaginationOptions
}

type (
	FilterOptions struct {
		AuthorID int
		LabelID  int
	}
	PaginationOptions struct {
		Offset int
		Limit  int
	}
)
