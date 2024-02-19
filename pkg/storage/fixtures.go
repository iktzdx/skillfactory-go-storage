package storage

import (
	"math"
	"math/rand/v2"
	"time"
)

const randomFactor float64 = 100000

func FixtureTask_1() Task {
	return Task{
		ID:         int(randomFactor) + 1,
		Opened:     time.Now().UTC(),
		Closed:     time.Now().UTC().AddDate(0, 0, 2),
		AuthorID:   0,
		AssignedID: 0,
		Title:      "Fixture #1",
		Content:    "Describe a tech-nature utopia.",
	}
}

func FixtureTask_2() Task {
	return Task{
		ID:         int(randomFactor) + 2,
		Opened:     time.Now().UTC(),
		Closed:     time.Now().UTC().AddDate(0, 0, 5),
		AuthorID:   0,
		AssignedID: 1,
		Title:      "Fixture #2",
		Content:    "Write a kitten's forest adventure.",
	}
}

func FixtureTask_3() Task {
	return Task{
		ID:         int(randomFactor) + 3,
		Opened:     time.Now().UTC(),
		Closed:     time.Time{},
		AuthorID:   1,
		AssignedID: 1,
		Title:      "Fixture #3",
		Content:    "Script AI ethics short film.",
	}
}

func FixtureTasks() Tasks {
	return Tasks{FixtureTask_1(), FixtureTask_2(), FixtureTask_3()}
}

func genRandTaskID(f float64) int {
	return int(math.Round((rand.Float64() * f) + f))
}
