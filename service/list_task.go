package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/tetran/go-web-app-example/auth"
	"github.com/tetran/go-web-app-example/entity"
)

type ListTask struct {
	DB   *sqlx.DB
	Repo TaskLister
}

func (lt *ListTask) ListTasks(ctx context.Context) (entity.Tasks, error) {
	id, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get user id")
	}
	ts, err := lt.Repo.ListTasks(ctx, lt.DB, id)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return ts, nil
}
