package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/tetran/go-web-app-example/entity"
	"github.com/tetran/go-web-app-example/store"
)

type AddTask struct {
	DB        *sqlx.DB
	Repo      *store.Repository
	Validator *validator.Validate
}

func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Set "go-playground/validator" style validation tag
	var b struct {
		Title string `json:"title" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	err := at.Validator.Struct(b)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	t := &entity.Task{
		Title:  b.Title,
		Status: entity.TaskStatusTodo,
	}
	id, err := at.Repo.AddTask(ctx, at.DB, t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	resp := struct {
		ID int `json:"id"`
	}{ID: id}
	RespondJSON(ctx, w, resp, http.StatusCreated)
}
