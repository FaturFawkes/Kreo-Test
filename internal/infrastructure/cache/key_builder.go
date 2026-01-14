package cache

import (
	"fmt"
)

// KeyBuilder provides methods to build Redis cache keys with consistent namespacing
type KeyBuilder struct {
	namespace string
}

// NewKeyBuilder creates a new KeyBuilder with the given namespace
func NewKeyBuilder(namespace string) *KeyBuilder {
	return &KeyBuilder{
		namespace: namespace,
	}
}

// MarketList builds a key for market list cache
func (kb *KeyBuilder) MarketList(category string) string {
	return fmt.Sprintf("%s:markets:list:%s", kb.namespace, category)
}

// MarketMetadata builds a key for market metadata cache
func (kb *KeyBuilder) MarketMetadata(ticker string) string {
	return fmt.Sprintf("%s:markets:metadata:%s", kb.namespace, ticker)
}

// MarketOrderBook builds a key for market order book cache
func (kb *KeyBuilder) MarketOrderBook(ticker string) string {
	return fmt.Sprintf("%s:markets:orderbook:%s", kb.namespace, ticker)
}

// MarketTrades builds a key for market trades cache
func (kb *KeyBuilder) MarketTrades(ticker string) string {
	return fmt.Sprintf("%s:markets:trades:%s", kb.namespace, ticker)
}

// CategoryOverview builds a key for category overview cache
func (kb *KeyBuilder) CategoryOverview(category string) string {
	return fmt.Sprintf("%s:categories:overview:%s", kb.namespace, category)
}

// CategoryList builds a key for category list cache
func (kb *KeyBuilder) CategoryList() string {
	return fmt.Sprintf("%s:categories:list", kb.namespace)
}

// RateLimitCounter builds a key for rate limit counter
func (kb *KeyBuilder) RateLimitCounter(identifier string, window string) string {
	return fmt.Sprintf("%s:ratelimit:counter:%s:%s", kb.namespace, identifier, window)
}

// RequestCoalescingLock builds a key for request coalescing lock
func (kb *KeyBuilder) RequestCoalescingLock(resource string) string {
	return fmt.Sprintf("%s:lock:coalesce:%s", kb.namespace, resource)
}

// HotMarkets builds a key for hot markets list
func (kb *KeyBuilder) HotMarkets() string {
	return fmt.Sprintf("%s:markets:hot", kb.namespace)
}
