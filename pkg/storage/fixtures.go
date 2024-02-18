package storage

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

const (
	fixtureTask1ID int = 1337
	fixtureTask2ID int = 7331
	fixtureTask3ID int = 1234

	openedAfterDays int = 5

	author1ID int = 1
	author2ID int = 2
)

func FixtureTask1() Task {
	return Task{
		ID:         fixtureTask1ID,
		Opened:     time.Now().UTC(),
		Closed:     time.Time{},
		AuthorID:   author1ID,
		AssignedID: author1ID,
		Title:      "Fixture Task #1",
		Content:    "This is a task #1 for tests.",
	}
}

func FixtureTask2() Task {
	return Task{
		ID:         fixtureTask2ID,
		Opened:     time.Now().UTC().AddDate(0, 0, openedAfterDays),
		Closed:     time.Time{},
		AuthorID:   author2ID,
		AssignedID: author2ID,
		Title:      "Fixture Task #2",
		Content:    "This is a task #2 for tests.",
	}
}

func FixtureTask3() Task {
	layout := "2006-Jan-02"
	timeOpened, _ := time.Parse(layout, "2023-Dec-31")

	return Task{
		ID:         fixtureTask3ID,
		Opened:     timeOpened,
		Closed:     time.Time{},
		AuthorID:   author1ID,
		AssignedID: author2ID,
		Title:      "Fixture Task #3",
		Content:    "This is a task #3 for tests.",
	}
}

func ApplyPostgreSQLFixtures(ctx context.Context, db *bun.DB, fxs Tasks) error {
	for _, fx := range fxs {
		fx := taskToRow(fx)

		_, err := db.NewInsert().Model(&fx).Exec(ctx)
		if err != nil {
			return errors.Wrap(err, "cannot apply fixture")
		}
	}

	return nil
}

func CleanPostgreSQLFixtures(ctx context.Context, db *bun.DB, fxs Tasks) error {
	for _, fx := range fxs {
		fx := taskToRow(fx)

		_, err := db.NewDelete().Model(&fx).WherePK().Exec(ctx)
		if err != nil {
			return errors.Wrap(err, "cannot delete fixture")
		}
	}

	return nil
}
