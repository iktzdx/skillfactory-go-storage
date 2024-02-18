package storage

import "context"

type Storage interface {
	Create(ctx context.Context, t Task) (int, error)
	GetByID(ctx context.Context, taskID int) (Task, error)
	List(ctx context.Context, opts SearchOptions) (Tasks, error)
	Update(ctx context.Context, t Task) (int, error)
	Delete(ctx context.Context, taskID int) error
}
