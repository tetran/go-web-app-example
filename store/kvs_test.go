package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/tetran/go-web-app-example/entity"
	"github.com/tetran/go-web-app-example/testutil"
)

func TestSave(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisForTest(t)

	sut := &KVS{cli}
	key := "TestKVS_Save"
	uid := entity.UserID(1234)
	ctx := context.Background()
	t.Cleanup(func() {
		cli.Del(ctx, key)
	})
	if err := sut.Save(ctx, key, uid); err != nil {
		t.Errorf("want no error, but has error: %v", err)
	}
}

func TestLoad(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisForTest(t)
	sut := &KVS{cli}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Load.ok"
		uid := entity.UserID(1234)
		ctx := context.Background()
		if err := cli.Set(ctx, key, int64(uid), 30*time.Minute).Err(); err != nil {
			t.Fatalf("failed to set value to redis: %v", err)
		}
		t.Cleanup(func() {
			cli.Del(ctx, key)
		})

		got, err := sut.Load(ctx, key)
		if err != nil {
			t.Fatalf("want no error, but has error: %v", err)
		}
		if got != uid {
			t.Errorf("want %d, but got %d", uid, got)
		}
	})

	t.Run("notFound", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Load.notFound"
		ctx := context.Background()
		got, err := sut.Load(ctx, key)
		if err == nil || !errors.Is(err, ErrNotFound) {
			t.Fatalf("want error, but no error: %v", got)
		}
	})
}
