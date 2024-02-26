package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/tetran/go-web-app-example/auth"
	"github.com/tetran/go-web-app-example/clock"
	"github.com/tetran/go-web-app-example/config"
	"github.com/tetran/go-web-app-example/handler"
	"github.com/tetran/go-web-app-example/service"
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
	clocker := clock.RealClocker{}

	// setup resources
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	rcli, err := store.NewKVS(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	jwter, err := auth.NewJWTer(rcli, clocker)
	if err != nil {
		return nil, cleanup, err
	}

	// POST /tasks
	r := store.Repository{Clocker: clocker}
	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/tasks", at.ServeHTTP)

	// GET /tasks
	lt := &handler.ListTask{
		Service: &service.ListTask{DB: db, Repo: &r},
	}
	mux.Get("/tasks", lt.ServeHTTP)

	// POST /users
	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/users", ru.ServeHTTP)

	// POST /login
	l := &handler.Login{
		Service:   &service.Login{DB: db, Repo: &r, TokenGenerator: jwter},
		Validator: v,
	}
	mux.Post("/login", l.ServeHTTP)

	return mux, cleanup, nil
}
