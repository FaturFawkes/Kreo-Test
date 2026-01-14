package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// ErrLockNotAcquired is returned when a lock cannot be acquired
	ErrLockNotAcquired = errors.New("lock not acquired")
)

const (
	// lockTTL is the default time-to-live for locks
	lockTTL = 30 * time.Second
	// lockRetryDelay is the delay between lock acquisition attempts
	lockRetryDelay = 50 * time.Millisecond
	// lockMaxWait is the maximum time to wait for lock acquisition
	lockMaxWait = 5 * time.Second
)

// Coalescer provides request coalescing using Redis locks
// This prevents multiple concurrent requests from hitting the same resource
type Coalescer struct {
	client     *redis.Client
	keyBuilder *KeyBuilder
}

// NewCoalescer creates a new request coalescer
func NewCoalescer(client *redis.Client, keyBuilder *KeyBuilder) *Coalescer {
	return &Coalescer{
		client:     client,
		keyBuilder: keyBuilder,
	}
}

// AcquireLock tries to acquire a distributed lock for a resource
// Returns true if lock was acquired, false if another process holds it
func (c *Coalescer) AcquireLock(ctx context.Context, resource string) (bool, error) {
	lockKey := c.keyBuilder.RequestCoalescingLock(resource)

	result, err := c.client.SetNX(ctx, lockKey, "1", lockTTL).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return result, nil
}

// ReleaseLock releases a distributed lock for a resource
func (c *Coalescer) ReleaseLock(ctx context.Context, resource string) error {
	lockKey := c.keyBuilder.RequestCoalescingLock(resource)

	if err := c.client.Del(ctx, lockKey).Err(); err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	return nil
}

// WaitForLock waits for a lock to become available
// Returns when lock is released or context times out
func (c *Coalescer) WaitForLock(ctx context.Context, resource string) error {
	lockKey := c.keyBuilder.RequestCoalescingLock(resource)

	waitCtx, cancel := context.WithTimeout(ctx, lockMaxWait)
	defer cancel()

	ticker := time.NewTicker(lockRetryDelay)
	defer ticker.Stop()

	for {
		select {
		case <-waitCtx.Done():
			return waitCtx.Err()
		case <-ticker.C:
			exists, err := c.client.Exists(waitCtx, lockKey).Result()
			if err != nil {
				return fmt.Errorf("failed to check lock: %w", err)
			}

			if exists == 0 {
				return nil
			}
		}
	}
}

// WithLock executes a function while holding a lock
// If lock cannot be acquired, waits for it to be released before returning
func (c *Coalescer) WithLock(ctx context.Context, resource string, fn func() error) error {
	acquired, err := c.AcquireLock(ctx, resource)
	if err != nil {
		return err
	}

	if !acquired {
		if err := c.WaitForLock(ctx, resource); err != nil {
			return fmt.Errorf("wait for lock failed: %w", err)
		}
		return nil
	}

	defer func() {
		if err := c.ReleaseLock(context.Background(), resource); err != nil {
			fmt.Printf("Failed to release lock for %s: %v\n", resource, err)
		}
	}()

	return fn()
}
