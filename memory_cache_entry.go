package memorycache

import "time"

// memoryCacheEntry запись кэша
type memoryCacheEntry struct {
	// Value значение
	Value interface{}
	// Expiration время жизни
	Expiration time.Time
	// Durability стойкость к вытеснению
	Durability MemoryCacheEntryDurability
}
