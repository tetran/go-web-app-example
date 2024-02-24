package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/tetran/go-web-app-example/clock"
	"github.com/tetran/go-web-app-example/config"
	"github.com/tetran/go-web-app-example/handler"
	"github.com/tetran/go-web-app-example/store"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// explicitly ignoring the return values to avoid errors in the static analysis
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}

	r := store.Repository{Clocker: clock.RealClocker{}}
	at := &handler.AddTask{DB: db, Repo: &r, Validator: v}
	mux.Post("/tasks", at.ServeHTTP)

	lt := &handler.ListTask{DB: db, Repo: &r}
	mux.Get("/tasks", lt.ServeHTTP)

	return mux, cleanup, nil
}
