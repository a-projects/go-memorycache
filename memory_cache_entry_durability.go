package memorycache

// MemoryCacheEntryDurability стойкость вытеснения записей из кэша при достижении лимита,
// устаревшие записи имеют стойкость ниже чем Weak
type MemoryCacheEntryDurability int

const (
	Weak   MemoryCacheEntryDurability = -1 // самая низкая стойкость, вытесняется первым
	Normal MemoryCacheEntryDurability = 0  // вытесняется если нет записей с приоритетом Weak
	Strong MemoryCacheEntryDurability = 1  // вытесняется если нет записей с приоритетом Normal
)
