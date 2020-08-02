package kvstore

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// KVStore is a key value store with support to TTL per key
type KVStore struct {
	client *redis.Client
}

// Error represents errors that may be returned by the KVStore.
// They should always be checked using errors.Is since the
// error may be wrapped with more context.
type Error string

const (
	KeyNotFoundErr Error = "key not found"
)

// New creates a new KVStore connecting to the given redis addr
// and using the given password to authenticate.
func New(addr string, password string) *KVStore {
	const maxRetries = 5

	client := redis.NewClient(&redis.Options{
		Addr:       addr,
		MaxRetries: maxRetries,
		Password:   password,
	})

	return &KVStore{client: client}
}

// Put adds the given key and value to the storage with the given  ttl.
// After the elapsed period the key will automatically be removed from the storage.
// If the key already exists it's value is going to be replaced
// by the new one along with the new ttl.
func (kv *KVStore) Put(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	scmd := kv.client.Set(ctx, key, val, ttl)
	_, err := scmd.Result()
	return err
}

// Get retrieves the value associated with the given key.
// If the key does not exists returns a KeyNotFound error,
// or any non-nil error in case of other failures.
func (kv *KVStore) Get(ctx context.Context, key string) ([]byte, error) {
	scmd := kv.client.Get(ctx, key)
	val, err := scmd.Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, KeyNotFoundErr
		}
		return nil, err
	}
	return val, nil
}

// Error returns the string representation of the error
func (e Error) Error() string {
	return string(e)
}
