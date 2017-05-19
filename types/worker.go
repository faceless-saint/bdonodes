package types

import (
	"math"
	"time"
)

const (
	// BaseWorkTime is multiplied by the ratio of work rate to workload.
	BaseWorkTime = time.Minute * 10
	// BaseMoveTime is multiplied by the ratio of speed to distance.
	BaseMoveTime = time.Second
)

// TODO: values currently pulled out of thin air - find better defaults!
var (
	defaultHumanWorker = Worker{
		Work:    100,
		Speed:   50,
		Luck:    0,
		Stamina: 20,
	}
	defaultGoblinWorker = Worker{
		Work:    130,
		Speed:   80,
		Luck:    0,
		Stamina: 15,
	}
	defaultGiantWorker = Worker{
		Work:    80,
		Speed:   30,
		Luck:    0,
		Stamina: 40,
	}
)

// HumanWorker returns an average human worker with the given home.
func HumanWorker(home string) *Worker {
	w := defaultHumanWorker
	w.Home = home
	return &w
}

// GoblinWorker returns an average goblin worker with the given home.
func GoblinWorker(home string) *Worker {
	w := defaultGoblinWorker
	w.Home = home
	return &w
}

// GiantWorker returns an average giant worker with the given home.
func GiantWorker(home string) *Worker {
	w := defaultGiantWorker
	w.Home = home
	return &w
}

// A Worker is an NPC that can be assigned to work a node.
type Worker struct {
	Home    string  `json:"home"`
	Work    float64 `json:"work"`
	Speed   float64 `json:"speed"`
	Luck    float64 `json:"luck"`
	Stamina uint16  `json:"stamina"`
}

// Endurance is the maximum work time until stamina must be refreshed.
func (w *Worker) Endurance(work float64, dist uint16) time.Duration {
	return w.Time(work, dist) * time.Duration(w.Stamina)
}

// Time taken to complete a cycle at a node with the given properties.
func (w *Worker) Time(work float64, dist uint16) time.Duration {
	return w.WorkTime(work) + w.MoveTime(dist)
}

// WorkTime at a node with the given workload.
func (w *Worker) WorkTime(work float64) time.Duration {
	// (sic) This method is wrong, but matches in-game implementation
	return BaseWorkTime * parseDuration(work/w.Work+1.0)
	//-> return BaseWorkTime * floatDuration(math.Ceil(work/w.Work))
}

// MoveTime to a node with the given distance.
func (w *Worker) MoveTime(dist uint16) time.Duration {
	return BaseMoveTime * parseDuration(float64(dist*2)/w.Speed)
}

// ... this is overkill for the above use cases, but is good practice
func parseDuration(input float64) time.Duration {
	// Enforce int64 boundaries on the input value
	switch {
	case input > math.MaxInt64:
		return time.Duration(math.MaxInt64)
	case input < math.MinInt64:
		return time.Duration(math.MinInt64)
	}
	return time.Duration(input)
}
