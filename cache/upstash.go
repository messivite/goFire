package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/claywarren/upstash-go"
)

// UpstashCache implements Cache using Upstash Redis REST API.
type UpstashCache struct {
	client    upstash.Upstash
	keyPrefix string
}

// NewUpstashCache creates a new Upstash Redis cache.
// keyPrefix is prepended to all keys (e.g. "api:").
func NewUpstashCache(restURL, restToken, keyPrefix string) (*UpstashCache, error) {
	client, err := upstash.New(upstash.Options{
		Url:   restURL,
		Token: restToken,
	})
	if err != nil {
		return nil, fmt.Errorf("connecting to Upstash Redis: %w", err)
	}

	if keyPrefix == "" {
		keyPrefix = "cache:"
	}

	log.Println("Upstash Redis cache initialized")
	return &UpstashCache{client: client, keyPrefix: keyPrefix}, nil
}

// Get returns cached bytes for the key, or nil if not found.
func (u *UpstashCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := u.client.Get(ctx, u.keyPrefix+key)
	if err != nil {
		return nil, err
	}
	if val == "" {
		return nil, nil
	}
	return []byte(val), nil
}

// SetAsync saves the value in the background with the given TTL in seconds.
func (u *UpstashCache) SetAsync(ctx context.Context, key string, data []byte, ttlSeconds int) {
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		fullKey := u.keyPrefix + key
		err := u.client.SetEX(bgCtx, fullKey, ttlSeconds, string(data))
		if err != nil {
			log.Printf("Redis SetEX error for %s: %v", fullKey, err)
		}
	}()
}
