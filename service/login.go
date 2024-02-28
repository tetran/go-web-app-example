package service

import (
	"context"
	"fmt"

	"github.com/tetran/go-web-app-example/store"
)

type Login struct {
	DB             store.Queryer
	Repo           UserGetter
	TokenGenerator TokenGenerator
}

func (l *Login) Login(ctx context.Context, name, password string) (string, error) {
	u, err := l.Repo.GetUser(ctx, l.DB, name)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	if err := u.VerifyPassword(password); err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}
	jwt, err := l.TokenGenerator.GenerateToken(ctx, *u)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return string(jwt), nil
}
