package store

import (
	"context"

	"github.com/tetran/go-web-app-example/entity"
)

func (r *Repository) ListTasks(ctx context.Context, db Queryer) (entity.Tasks, error) {
	tasks := entity.Tasks{}

	sql := `SELECT id, title, status, created_at, updated_at FROM tasks;`
	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) AddTask(ctx context.Context, db Executer, t *entity.Task) (int, error) {
	t.CreatedAt = r.Clocker.Now()
	t.UpdatedAt = r.Clocker.Now()

	sql := `INSERT INTO tasks (title, status, created_at, updated_at) VALUES (:title, :status, :created_at, :updated_at);`
	res, err := db.NamedExecContext(ctx, sql, t)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	t.ID = entity.TaskID(id)
	return int(t.ID), nil
}
