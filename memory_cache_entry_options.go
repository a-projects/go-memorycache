package memorycache

import "time"

// MemoryCacheEntryOptions параметры записи
type MemoryCacheEntryOptions struct {
	// Expiration время жизни записи
	Expiration time.Time
	// Durability стойкость к вытеснению, по умолчанию Normal
	Durability MemoryCacheEntryDurability
}
