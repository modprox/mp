package problems

import (
	"sort"
	"sync"
	"time"

	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/loggy"
)

//go:generate go run github.com/gojuno/minimock/cmd/minimock -g -i Tracker -s _mock.go

type Tracker interface {
	Set(Problem)
	Problem(module coordinates.Module) (Problem, bool)
	Problems() []Problem
}

type Problem struct {
	Module  coordinates.Module `json:"module"`
	Time    time.Time          `json:"time"`
	Message string             `json:"message"`
}

func Create(mod coordinates.Module, err error) Problem {
	return Problem{
		Module:  mod,
		Time:    time.Now(),
		Message: err.Error(),
	}
}

type tracker struct {
	log      loggy.Logger
	lock     sync.RWMutex
	problems map[coordinates.Module]Problem
}

func New(name string) Tracker {
	return &tracker{
		log:      loggy.New("problems-" + name),
		problems: make(map[coordinates.Module]Problem),
	}
}

func (t *tracker) Set(problem Problem) {
	t.log.Tracef("setting problem for module %s", problem.Module)

	t.lock.Lock()
	defer t.lock.Unlock()

	t.problems[problem.Module] = problem
}

func (t *tracker) Problem(mod coordinates.Module) (Problem, bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	p, exists := t.problems[mod]
	return p, exists
}

func (t *tracker) Problems() []Problem {
	t.lock.RLock()
	defer t.lock.RUnlock()

	problems := make([]Problem, 0, len(t.problems))
	for _, p := range t.problems {
		problems = append(problems, p)
	}

	sort.Sort(byName(problems))
	return problems
}

type byName []Problem

func (p byName) Len() int      { return len(p) }
func (p byName) Swap(x, y int) { p[x], p[y] = p[y], p[x] }
func (p byName) Less(x, y int) bool {
	modX, modY := p[x], p[y]

	if modX.Module.Source < modY.Module.Source {
		return true
	} else if modX.Module.Source > modY.Module.Source {
		return false
	}

	// don't bother parsing the tags, nobody cares
	// we just need something deterministic
	if modX.Time.Before(modY.Time) {
		return true
	}
	return false
}
