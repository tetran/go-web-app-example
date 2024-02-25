package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/tetran/go-web-app-example/clock"
	"github.com/tetran/go-web-app-example/entity"
)

func TestRegisterUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	c := clock.FixedClocker{}
	var wantId int64 = 20
	okUser := &entity.User{
		Name: "Pong", Password: "passwd", Role: "admin", CreatedAt: c.Now(), UpdatedAt: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectExec(`INSERT INTO users \(name, password, role, created_at, updated_at\) VALUES \(\?, \?, \?, \?, \?\);`).
		WithArgs(okUser.Name, okUser.Password, okUser.Role, okUser.CreatedAt, okUser.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(wantId, 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if _, err := r.RegisterUser(ctx, xdb, okUser); err != nil {
		t.Errorf("want no error, but has error: %v", err)
	}
}
