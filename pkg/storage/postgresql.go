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
		return -1, errors.Wrap(err, "last inserted task")
	}

	return int(affected), nil
}

func (r *postgreSQLRepo) BulkCreate(ctx context.Context, ts Tasks) (int, error) {
	tasksRows := tasksToRows(ts)

	res, err := r.db.NewInsert().Model(&tasksRows).Exec(ctx)
	if err != nil {
		return -1, errors.Wrap(err, "insert multiple tasks")
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return -1, errors.Wrap(err, "get affected rows")
	}

	return int(affected), nil
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

	/*
	   SELECT t.*
	   FROM "tasks" AS "t"
	   JOIN tasks_labels AS tls ON tls.task_id = t.id
	   JOIN labels AS l on l.id = tls.label_id
	   WHERE (t.author_id = 1)
	   OR (tls.label_id = 0)
	   ORDER BY "t"."id" ASC
	*/

	query := r.db.NewSelect().
		Model((*row)(nil)).
		ColumnExpr("?", bun.Safe("t.*")).
		Join("JOIN tasks_labels AS tls ON tls.task_id = t.id").
		Join("JOIN labels AS l on l.id = tls.label_id").
		Where("t.author_id = ?", opts.AuthorID).
		WhereOr("tls.label_id = ?", opts.LabelID).
		Order("t.id ASC")

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

	if len(tasksRows) == 0 {
		err := errors.New("error while scan list")

		return Tasks{}, errors.Wrapf(err, "q: %v\n", query.String())
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

	return int(affected), nil
}

func (r *postgreSQLRepo) Delete(ctx context.Context, taskID int) (int, error) {
	res, err := r.db.NewDelete().Model((*row)(nil)).Where("id = ?", taskID).Exec(ctx)
	if err != nil {
		return -1, errors.Wrap(err, "delete a task")
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return -1, errors.Wrap(err, "get affected rows")
	}

	return int(affected), nil
}
