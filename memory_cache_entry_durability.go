package memorycache

// MemoryCacheEntryDurability resistance to evicting entries from cache when limit is reached
type MemoryCacheEntryDurability int

const (
	Weak   MemoryCacheEntryDurability = -1 // lowest resistance, evicted first
	Normal MemoryCacheEntryDurability = 0  // evicted if there are no records with priority Weak
	Strong MemoryCacheEntryDurability = 1  // evicted if there are no records with priority Normal
)
