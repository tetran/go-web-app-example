package service

//go:generate go run github.com/matryer/moq -out moq_test.go . TaskAdder TaskLister UserRegisterer UserGetter

import (
	"context"

	"github.com/tetran/go-web-app-example/entity"
	"github.com/tetran/go-web-app-example/store"
)

type TaskAdder interface {
	AddTask(ctx context.Context, db store.Executer, t *entity.Task) (int, error)
}

type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer) (entity.Tasks, error)
}

type UserRegisterer interface {
	RegisterUser(ctx context.Context, db store.Executer, u *entity.User) (int64, error)
}

type UserGetter interface {
	GetUser(ctx context.Context, db store.Queryer, name string) (*entity.User, error)
}

type TokenGenerator interface {
	GenerateToken(ctx context.Context, u entity.User) ([]byte, error)
}
