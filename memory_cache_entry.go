package memorycache

import "time"

// memoryCacheEntry
type memoryCacheEntry struct {
	// Value entry value
	Value interface{}
	// Expirationentry lifetime
	Expiration time.Time
	// Durability eviction resistance
	Durability MemoryCacheEntryDurability
}
