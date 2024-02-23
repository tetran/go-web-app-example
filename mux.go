package main

import "net/http"

func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// Ignoring the return values to avoid errors in the static analysis
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	return mux
}
