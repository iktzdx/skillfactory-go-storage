package storage

import (
	"context"
	"database/sql"
	"math"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type PostgreSQLTestSuite struct {
	suite.Suite
	repo Storage
	db   *bun.DB
}

func (s *PostgreSQLTestSuite) applyFixtures(ctx context.Context, tasks Tasks) error {
	rows := tasksToRows(tasks)

	_, err := s.db.NewInsert().Model(&rows).Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "apply fixtures")
	}

	return nil
}

func (s *PostgreSQLTestSuite) labelFixture(ctx context.Context, taskID int) error {
	_, err := s.db.NewRaw("INSERT INTO tasks_labels (task_id, label_id) VALUES (?, 0)", taskID).Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "label fixtures")
	}

	return nil
}

func (s *PostgreSQLTestSuite) flushFixtures(ctx context.Context, minID int) error {
	_, err := s.db.NewDelete().
		Model((*row)(nil)).
		Where("id >= ?", minID).
		Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "flush fixtures")
	}

	return nil
}

func (s *PostgreSQLTestSuite) compareTasks(expected, actual Task) {
	s.Require().Equal(expected.ID, actual.ID, "compare id")
	s.Require().Equal(expected.AuthorID, actual.AuthorID, "compare author id")
	s.Require().Equal(expected.AssignedID, actual.AssignedID, "compare assigned id")
	s.Require().Equal(expected.Title, actual.Title, "compare title")
	s.Require().Equal(expected.Content, actual.Content, "compare content")

	err := isDateEqual(expected.Opened, actual.Opened)
	s.Require().NoError(err, "compare opened date")

	err = isDateEqual(expected.Closed, actual.Closed)
	s.Require().NoError(err, "compare closed date")
}

func isDateEqual(expected, actual time.Time) error {
	if !expected.Equal(actual) && math.Abs(float64(expected.Sub(actual).Nanoseconds())) > 1e9 {
		return errors.New("dates not equal")
	}

	return nil
}

func TestPostgreSQLTestSuite(t *testing.T) {
	suite.Run(t, new(PostgreSQLTestSuite))
}

func (s *PostgreSQLTestSuite) SetupSuite() {
	ctx := context.Background()

	dsn := "postgres://postgres:example@localhost:5432/postgres?sslmode=disable"

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	s.db = bun.NewDB(sqldb, pgdialect.New())

	s.repo = NewRepoPostgreSQL(s.db)

	err := s.applyFixtures(ctx, FixtureTasks())
	s.Require().NoError(err)
}

func (s *PostgreSQLTestSuite) TearDownSuite() {
	ctx := context.Background()

	err := s.flushFixtures(ctx, int(randomFactor))
	s.Require().NoError(err)

	err = s.db.Close()
	s.Require().NoError(err)
}

func (s *PostgreSQLTestSuite) TestGetCreatedTaskByID() {
	s.Run("Task exists", func() {
		ctx := context.Background()

		want := Task{
			ID:         genRandTaskID(randomFactor),
			Opened:     time.Now().UTC(),
			Closed:     time.Time{},
			AuthorID:   1,
			AssignedID: 1,
			Title:      "Test task #1",
			Content:    "AR game blending reality with interactive virtual elements.",
		}

		affected, err := s.repo.Create(ctx, want)
		s.Require().NoError(err)
		s.Require().Greater(affected, 0)

		got, err := s.repo.GetByID(ctx, want.ID)
		s.Require().NoError(err)

		s.compareTasks(want, got)
	})

	s.Run("Task not founded", func() {
		ctx := context.Background()
		_, err := s.repo.GetByID(ctx, genRandTaskID(randomFactor))
		s.Require().ErrorIs(err, sql.ErrNoRows)
	})
}

func (s *PostgreSQLTestSuite) TestListCreatedTasks() {
	ctx := context.Background()

	want := Tasks{
		{
			ID:         genRandTaskID(randomFactor),
			Opened:     time.Now().UTC(),
			Closed:     time.Time{},
			AuthorID:   1,
			AssignedID: 0,
			Title:      "Test List Created Tasks #1",
			Content:    "Write a short summary of the book you're currently reading.",
		},
		{
			ID:         genRandTaskID(randomFactor),
			Opened:     time.Now().UTC(),
			Closed:     time.Time{},
			AuthorID:   1,
			AssignedID: 1,
			Title:      "Test List Created Tasks #2",
			Content:    "Describe your favorite recipe in 5 sentences or less.",
		},
		{
			ID:         genRandTaskID(randomFactor),
			Opened:     time.Now().UTC(),
			Closed:     time.Time{},
			AuthorID:   0,
			AssignedID: 1,
			Title:      "Test List Created Tasks #3",
			Content:    "Explain the concept of \"machine learning\" to a 5-year-old.",
		},
	}

	affected, err := s.repo.BulkCreate(ctx, want)
	s.Require().NoError(err)
	s.Require().Equal(len(want), affected)

	for i := 0; i < len(want); i++ {
		err := s.labelFixture(ctx, want[i].ID)
		s.Require().NoError(err)
	}

	//nolint:exhaustruct
	listedTasks, err := s.repo.List(ctx, SearchOptions{
		FilterOptions: FilterOptions{AuthorID: 1},
	})

	s.Require().NoError(err)
	s.Require().NotEmpty(listedTasks)

	const expectedAmount int = 2
	for i := 0; i < expectedAmount; i++ {
		s.compareTasks(want[i], listedTasks[i])
	}
}
