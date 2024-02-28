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
	"github.com/tetran/go-web-app-example/testutil/fixture"
)

func prepareUser(ctx context.Context, t *testing.T, db Executer) entity.UserID {
	t.Helper()
	u := fixture.User(nil)
	result, err := db.NamedExecContext(ctx,
		"INSERT INTO users (name, password, role, created_at, updated_at) VALUES (:name, :password, :role, :created_at, :updated_at);",
		u,
	)
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("failed to get last insert id: %v", err)
	}
	return entity.UserID(id)
}

func prepareTasks(ctx context.Context, t *testing.T, con Executer) (entity.UserID, entity.Tasks) {
	t.Helper()

	if _, err := con.ExecContext(ctx, "TRUNCATE TABLE tasks;"); err != nil {
		t.Fatalf("failed to truncate table: %v", err)
	}

	userId := prepareUser(ctx, t, con)
	otherUserId := prepareUser(ctx, t, con)

	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			UserID: userId, Title: "wants task 1", Status: "todo", CreatedAt: c.Now(), UpdatedAt: c.Now(),
		},
		{
			UserID: userId, Title: "wants task 2", Status: "done", CreatedAt: c.Now(), UpdatedAt: c.Now(),
		},
		{
			UserID: userId, Title: "wants task 3", Status: "doing", CreatedAt: c.Now(), UpdatedAt: c.Now(),
		},
	}
	not_wants := entity.Tasks{
		{
			UserID: otherUserId, Title: "not wants task 1", Status: "todo", CreatedAt: c.Now(), UpdatedAt: c.Now(),
		},
	}

	inserts := append(wants, not_wants...)
	result, err := con.NamedExecContext(ctx,
		"INSERT INTO tasks (user_id, title, status, created_at, updated_at) VALUES (:user_id, :title, :status, :created_at, :updated_at);",
		inserts,
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

	return userId, wants
}

// Use real database for list tasks
func TestListTasks(t *testing.T) {
	ctx := context.Background()

	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}

	uid, wants := prepareTasks(ctx, t, tx)
	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx, uid)
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
		UserID: 1, Title: "test task", Status: "todo", CreatedAt: c.Now(), UpdatedAt: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectExec(`INSERT INTO tasks \(user_id, title, status, created_at, updated_at\) VALUES \(\?, \?, \?, \?, \?\);`).
		WithArgs(okTask.UserID, okTask.Title, okTask.Status, okTask.CreatedAt, okTask.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(wantId, 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if _, err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but has error: %v", err)
	}
}
