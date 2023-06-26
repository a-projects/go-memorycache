package memorycache

import "time"

// MemoryCacheOptions cache options
type MemoryCacheOptions struct {
	// CleanupInterval interval for clearing cache from obsolete entries, if 0, then it does not start
	CleanupInterval time.Duration

	// LimitEntries limit of entries in cache, upon reaching which eviction begins, if 0, then unlimited
	LimitEntries int

	// StoreFile file to restore cache when application is restarted, if empty, it does not restore
	StoreFile string
}
