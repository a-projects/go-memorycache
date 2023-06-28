package memorycache

import "time"

// MemoryCacheOptions cache options
type MemoryCacheOptions struct {
	// CleanupInterval interval for clearing cache from obsolete entries, if not set, then it does not start
	CleanupInterval time.Duration

	// LimitEntries limit of entries in cache, upon reaching which eviction begins, if not set, then unlimited
	LimitEntries int

	// FileName file to restore cache when application is restarted, if not set, it does not restore
	FileName string
}
