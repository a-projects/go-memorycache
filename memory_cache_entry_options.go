package memorycache

import "time"

// MemoryCacheEntryOptions
type MemoryCacheEntryOptions struct {
	// Expiration entry lifetime
	Expiration time.Time
	// Durability eviction resistance, default Normal
	Durability MemoryCacheEntryDurability
}
