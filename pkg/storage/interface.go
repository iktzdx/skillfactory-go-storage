package storage

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, task Task) (int, error)
	BulkCreate(ctx context.Context, tasks Tasks) (int, error)
	GetByID(ctx context.Context, taskID int) (Task, error)
	List(ctx context.Context, opts SearchOptions) (Tasks, error)
	Update(ctx context.Context, task Task) (int, error)
	Delete(ctx context.Context, taskID int) (int, error)
}
