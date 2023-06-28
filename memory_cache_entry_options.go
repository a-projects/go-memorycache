package memorycache

import "time"

// MemoryCacheEntryOptions
type MemoryCacheEntryOptions struct {
	// Lifetime entry lifetime, if not set, then 0
	Lifetime time.Duration
	// Durability eviction resistance, if not set, then Normal
	Durability MemoryCacheEntryDurability
}
