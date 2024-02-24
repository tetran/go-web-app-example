package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"github.com/tetran/go-web-app-example/clock"
	"github.com/tetran/go-web-app-example/entity"
	"github.com/tetran/go-web-app-example/testutil"
)

func prepareTasks(ctx context.Context, t *testing.T, con Executer) entity.Tasks {
	t.Helper()

	if _, err := con.ExecContext(ctx, "TRUNCATE TABLE tasks;"); err != nil {
		t.Fatalf("failed to truncate table: %v", err)
	}

	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			Title: "wants task 1", Status: "todo", CreatedAt: c.Now(), UpdatedAt: c.Now(),
		},
		{
			Title: "wants task 2", Status: "done", CreatedAt: c.Now(), UpdatedAt: c.Now(),
		},
		{
			Title: "wants task 3", Status: "doing", CreatedAt: c.Now(), UpdatedAt: c.Now(),
		},
	}

	result, err := con.NamedExecContext(ctx,
		"INSERT INTO tasks (title, status, created_at, updated_at) VALUES (:title, :status, :created_at, :updated_at);",
		wants,
	)
	if err != nil {
		t.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)

	return wants
}

// Use real database for list tasks
func TestListTasks(t *testing.T) {
	ctx := context.Background()

	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}

	wants := prepareTasks(ctx, t, tx)
	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) > 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

// Use mock database for add task
func TestAddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	c := clock.FixedClocker{}
	var wantId int64 = 20
	okTask := &entity.Task{
		Title: "test task", Status: "todo", CreatedAt: c.Now(), UpdatedAt: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectExec(`INSERT INTO tasks \(title, status, created_at, updated_at\) VALUES \(\?, \?, \?, \?\);`).
		WithArgs(okTask.Title, okTask.Status, okTask.CreatedAt, okTask.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(wantId, 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if _, err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but has error: %v", err)
	}
}
