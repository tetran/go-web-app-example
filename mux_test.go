package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/tetran/go-web-app-example/config"
)

func TestNewMux(t *testing.T) {
	// t.Skip("TODO: ")

	cfg, err := config.New()
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	cfg.DBHost = "127.0.0.1"
	cfg.DBPort = 13306
	cfg.RedisHost = "127.0.0.1"
	cfg.RedisPort = 16379
	if _, defined := os.LookupEnv("CI"); defined {
		cfg.DBPort = 3306
		cfg.RedisPort = 6379
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	sut, _, err := NewMux(context.Background(), cfg)
	if err != nil {
		t.Fatalf("failed to create mux: %v", err)
	}
	sut.ServeHTTP(w, r)
	resp := w.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	want := `{"status": "ok"}`
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
}
