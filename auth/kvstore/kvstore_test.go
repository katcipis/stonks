package kvstore_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/katcipis/stonks/auth/kvstore"
)

func TestKVStorePutKeyWithTTL(t *testing.T) {
	const password = "test-kvstore-db-pass"

	m := newTestRedis(t, password)
	defer m.Close()

	s := kvstore.New(m.Addr(), password)

	const key = "testkey"
	const ttl = 10 * time.Second

	val := []byte("testval")
	ctx := context.Background()

	err := s.Put(ctx, key, val, ttl)
	assertNoErr(t, err)

	gotVal, err := s.Get(ctx, key)
	assertNoErr(t, err)

	want := string(val)
	got := string(gotVal)

	if got != want {
		t.Fatalf("retrieving data got %q want %q", got, want)
	}

	gotTTL := m.TTL(key)
	if gotTTL != ttl {
		t.Fatalf("got TTL %d want %d", gotTTL, ttl)
	}
}

func TestKVStoreKeyDoesNotExist(t *testing.T) {
	const password = "test-kvstore-db-pass-again"

	m := newTestRedis(t, password)
	defer m.Close()

	s := kvstore.New(m.Addr(), password)
	_, err := s.Get(context.Background(), "key")
	if !errors.Is(kvstore.KeyNotFoundErr, err) {
		t.Fatalf("got err[%v] want[%v]", err, kvstore.KeyNotFoundErr)
	}
}

func newTestRedis(t *testing.T, password string) *miniredis.Miniredis {
	t.Helper()

	m, err := miniredis.Run()
	assertNoErr(t, err)
	m.RequireAuth(password)
	return m
}

func assertNoErr(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatal(err)
	}
}
