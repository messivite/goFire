package cache

import "context"

// Cache defines a minimal key-value cache interface for use with Upstash Redis.
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	SetAsync(ctx context.Context, key string, data []byte, ttlSeconds int)
}
