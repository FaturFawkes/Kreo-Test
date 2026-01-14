package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrCoalescingTimeout is returned when waiting for coalesced request times out
	ErrCoalescingTimeout = errors.New("coalescing timeout")
)

// CoalescingLock defines the interface for distributed locking
type CoalescingLock interface {
	AcquireLock(ctx context.Context, resource string) (bool, error)
	ReleaseLock(ctx context.Context, resource string) error
	WaitForLock(ctx context.Context, resource string) error
}

// RequestCoalescer prevents thundering herd by coalescing concurrent requests
type RequestCoalescer struct {
	lock     CoalescingLock
	cache    RequestCache
	mu       sync.RWMutex
	inflight map[string]*inflightRequest
	maxWait  time.Duration
}

// RequestCache defines the interface for caching request results
type RequestCache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

// inflightRequest tracks a request that is currently being processed
type inflightRequest struct {
	done   chan struct{}
	result interface{}
	err    error
}

// NewRequestCoalescer creates a new request coalescer
func NewRequestCoalescer(lock CoalescingLock, cache RequestCache) *RequestCoalescer {
	return &RequestCoalescer{
		lock:     lock,
		cache:    cache,
		inflight: make(map[string]*inflightRequest),
		maxWait:  5 * time.Second,
	}
}

// Execute executes a function with request coalescing
// If multiple concurrent requests arrive for the same key, only one executes the function
func (rc *RequestCoalescer) Execute(
	ctx context.Context,
	key string,
	fn func(ctx context.Context) (interface{}, error),
) (interface{}, error) {
	if rc.cache != nil {
		result, err := rc.cache.Get(ctx, key)
		if err == nil && result != nil {
			return result, nil
		}
	}

	rc.mu.RLock()
	if req, exists := rc.inflight[key]; exists {
		rc.mu.RUnlock()
		select {
		case <-req.done:
			return req.result, req.err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	rc.mu.RUnlock()

	acquired, err := rc.lock.AcquireLock(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("lock acquisition failed: %w", err)
	}

	if !acquired {
		if err := rc.lock.WaitForLock(ctx, key); err != nil {
			return nil, ErrCoalescingTimeout
		}

		if rc.cache != nil {
			result, err := rc.cache.Get(ctx, key)
			if err == nil && result != nil {
				return result, nil
			}
		}

		return fn(ctx)
	}

	req := &inflightRequest{
		done: make(chan struct{}),
	}

	rc.mu.Lock()
	rc.inflight[key] = req
	rc.mu.Unlock()

	defer func() {
		rc.mu.Lock()
		delete(rc.inflight, key)
		rc.mu.Unlock()

		close(req.done)

		_ = rc.lock.ReleaseLock(context.Background(), key)
	}()

	result, execErr := fn(ctx)

	req.result = result
	req.err = execErr

	if execErr == nil && rc.cache != nil && result != nil {
		_ = rc.cache.Set(ctx, key, result, 0)
	}

	return result, execErr
}
