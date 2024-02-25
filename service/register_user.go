package service

import (
	"context"
	"fmt"

	"github.com/tetran/go-web-app-example/entity"
	"github.com/tetran/go-web-app-example/store"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	DB   store.Executer
	Repo UserRegisterer
}

func (ru *RegisterUser) RegisterUser(ctx context.Context, name, password, role string) (*entity.User, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	u := &entity.User{
		Name:     name,
		Password: string(pw),
		Role:     role,
	}

	if _, err := ru.Repo.RegisterUser(ctx, ru.DB, u); err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return u, nil
}
