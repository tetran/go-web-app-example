package auth

import (
	"bytes"
	"context"
	"testing"

	"github.com/tetran/go-web-app-example/clock"
	"github.com/tetran/go-web-app-example/entity"
	"github.com/tetran/go-web-app-example/testutil/fixture"
)

func TestEmbed(t *testing.T) {
	want := []byte("-----BEGIN PUBLIC KEY-----")
	if !bytes.Contains(rawPubKey, want) {
		t.Errorf("want %s, but got %s", want, rawPubKey)
	}

	want = []byte("-----BEGIN PRIVATE KEY-----")
	if !bytes.Contains(rawPrivKey, want) {
		t.Errorf("want %s, but got %s", want, rawPrivKey)
	}
}

func TestJWTer_GenerateToken(t *testing.T) {
	wantID := entity.UserID(1234)
	u := fixture.User(&entity.User{ID: wantID})
	moq := &StoreMock{}
	moq.SaveFunc = func(ctx context.Context, key string, userID entity.UserID) error {
		if userID != wantID {
			t.Errorf("want %d, but got %d", wantID, userID)
		}
		return nil
	}

	sut, err := NewJWTer(moq, clock.RealClocker{})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	got, err := sut.GenerateToken(ctx, *u)
	if err != nil {
		t.Fatalf("want no error, but got %v", err)
	}
	if len(got) == 0 {
		t.Error("want token, but got empty")
	}
}
