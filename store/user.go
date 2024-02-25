package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/tetran/go-web-app-example/entity"
)

func (r *Repository) RegisterUser(ctx context.Context, db Executer, u *entity.User) (int64, error) {
	u.CreatedAt = r.Clocker.Now()
	u.UpdatedAt = r.Clocker.Now()

	const q = `INSERT INTO users (name, password, role, created_at, updated_at) VALUES (:name, :password, :role, :created_at, :updated_at);`
	res, err := db.NamedExecContext(ctx, q, u)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return 0, fmt.Errorf("%s: name is already taken", ErrDuplicatedEntry)
		}
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	u.ID = entity.UserID(id)
	return id, nil
}
