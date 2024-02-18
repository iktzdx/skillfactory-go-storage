package storage

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type postgreSQLRepo struct {
	db *bun.DB
}

func NewRepoPostgreSQL(db *bun.DB) *postgreSQLRepo {
	return &postgreSQLRepo{db}
}

func (r *postgreSQLRepo) Create(ctx context.Context, t Task) (int, error) {
	taskRow := taskToRow(t)

	res, err := r.db.NewInsert().Model(&taskRow).Exec(ctx)
	if err != nil {
		return -1, errors.Wrap(err, "insert a task")
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return -1, errors.Wrap(err, "get affected rows")
	}

	if affected == 0 {
		return -1, errors.New("number of rows affected must be greater than 0")
	}

	taskID, err := res.LastInsertId()
	if err != nil {
		return -1, errors.Wrap(err, "get last inserted ID")
	}

	return int(taskID), nil
}

func (r *postgreSQLRepo) GetByID(ctx context.Context, taskID int) (Task, error) {
	var taskRow row

	err := r.db.NewSelect().Model((*row)(nil)).Where("id = ?", taskID).Scan(ctx, &taskRow)
	if err != nil {
		return Task{}, errors.Wrap(err, "get a task by id")
	}

	return rowToTask(taskRow), nil
}

func (r *postgreSQLRepo) List(ctx context.Context, opts SearchOptions) (Tasks, error) {
	var tasksRows rows

	query := r.db.NewSelect().
		Model((*row)(nil)).
		ColumnExpr("tasks.*").
		Join("JOIN tasks_labels AS tls ON tls.task_id = tasks.id").
		Join("JOIN labels AS l on l.id = tls.label_id").
		Order("id ASC")

	query = query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
		authorID := opts.AuthorID
		labelID := opts.LabelID

		q = q.Where("? = 0", authorID).
			WhereOr("author_id = ?", authorID)

		q = q.Where("? = 0", labelID).
			WhereOr("l.id = ?", labelID)

		return q
	})

	if opts.Limit != 0 {
		query = query.Limit(opts.Limit)
	}

	if opts.Offset != 0 {
		query = query.Offset(opts.Offset)
	}

	err := query.Scan(ctx, &tasksRows)
	if err != nil {
		return Tasks{}, errors.Wrap(err, "list tasks")
	}

	return rowsToTasks(tasksRows), nil
}

func (r *postgreSQLRepo) Update(ctx context.Context, t Task) (int, error) {
	taskRow := taskToRow(t)

	res, err := r.db.NewUpdate().Model(&taskRow).WherePK().Exec(ctx)
	if err != nil {
		return -1, errors.Wrap(err, "update a task")
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return -1, errors.Wrap(err, "get affected rows")
	}

	if affected == 0 {
		return -1, errors.New("number of rows affected must be greater than 0")
	}

	taskID, err := res.LastInsertId()
	if err != nil {
		return -1, errors.Wrap(err, "get last updated ID")
	}

	return int(taskID), nil
}

func (r *postgreSQLRepo) Delete(ctx context.Context, taskID int) error {
	res, err := r.db.NewDelete().Model((*row)(nil)).Where("id = ?", taskID).Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "delete a task")
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get affected rows")
	}

	if affected == 0 {
		return errors.New("number of rows affected must be greater than 0")
	}

	return nil
}
