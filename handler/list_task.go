package handler

import (
	"net/http"

	"github.com/tetran/go-web-app-example/entity"
)

type ListTask struct {
	Service ListTasksService
}

type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks, err := lt.Service.ListTasks(ctx)
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
