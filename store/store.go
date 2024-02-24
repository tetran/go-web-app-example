package store

import (
	"errors"

	"github.com/tetran/go-web-app-example/entity"
)

type TaskStore struct {
	LastID entity.TaskID
	Tasks  map[entity.TaskID]*entity.Task
}

var (
	Tasks       = &TaskStore{Tasks: map[entity.TaskID]*entity.Task{}}
	ErrNotFound = errors.New("task not found")
)

func (ts *TaskStore) Add(t *entity.Task) (int, error) {
	ts.LastID++
	t.ID = ts.LastID
	ts.Tasks[t.ID] = t
	return int(t.ID), nil
}

func (ts *TaskStore) All() entity.Tasks {
	tasks := make(entity.Tasks, 0, len(ts.Tasks))
	for _, t := range ts.Tasks {
		tasks = append(tasks, t)
	}
	return tasks
}
