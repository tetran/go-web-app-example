package handler

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/tetran/go-web-app-example/entity"
	"github.com/tetran/go-web-app-example/store"
)

type ListTask struct {
	DB   *sqlx.DB
	Repo *store.Repository
}

type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks, err := lt.Repo.ListTasks(ctx, lt.DB)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	ts := make([]task, 0, len(tasks))
	for _, t := range tasks {
		ts = append(ts, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		})
	}

	RespondJSON(ctx, w, ts, http.StatusOK)
}
