package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func AssertJSON(t *testing.T, want, got []byte) {
	t.Helper()

	var jw, jq any
	if err := json.Unmarshal(want, &jw); err != nil {
		t.Fatalf("failed to unmarshal want %q: %v", want, err)
	}
	if err := json.Unmarshal(got, &jq); err != nil {
		t.Fatalf("failed to unmarshal got %q: %v", got, err)
	}
	if diff := cmp.Diff(jq, jw); diff != "" {
		t.Errorf("mismatch (-got +want):\n%s", diff)
	}
}

func AssertResponse(t *testing.T, got *http.Response, status int, body []byte) {
	t.Helper()
	t.Cleanup(func() { _ = got.Body.Close() })

	gb, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatal(err)
	}
	if got.StatusCode != status {
		t.Fatalf("want status code %d, but got %d.\nbody: %q", status, got.StatusCode, gb)
	}
	if len(gb) == 0 && len(body) == 0 {
		return
	}

	AssertJSON(t, body, gb)
}

func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %q: %v", path, err)
	}
	return b
}
