package storage

import (
	"time"

	"github.com/uptrace/bun"
)

type rows []row

type row struct {
	bun.BaseModel `bun:"table:tasks"`

	ID         int       `bun:"id,pk,autoincrement"`
	Opened     time.Time `bun:"opened,nullzero,notnull,default:current_timestamp"`
	Closed     time.Time `bun:"closed,nullzero"`
	AuthorID   int       `bun:"author_id"`
	AssignedID int       `bun:"assigned_id"`
	Title      string    `bun:"title,notnull"`
	Content    string    `bun:"content"`
}

func taskToRow(t Task) row {
	return row{ //nolint:exhaustruct
		ID:         t.ID,
		Opened:     t.Opened,
		Closed:     t.Closed,
		AuthorID:   t.AuthorID,
		AssignedID: t.AssignedID,
		Title:      t.Title,
		Content:    t.Content,
	}
}

func rowToTask(r row) Task {
	return Task{
		ID:         r.ID,
		Opened:     r.Opened,
		Closed:     r.Closed,
		AuthorID:   r.AuthorID,
		AssignedID: r.AssignedID,
		Title:      r.Title,
		Content:    r.Content,
	}
}

func rowsToTasks(rs rows) Tasks {
	tasks := make(Tasks, len(rs))

	for idx, r := range rs {
		tasks[idx] = rowToTask(r)
	}

	return tasks
}
