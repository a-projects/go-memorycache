package memorycache

import "time"

// MemoryCacheOptions параметры кэша
type MemoryCacheOptions struct {
	// CleanupInterval интервал очиски кэша от устаревших записей, если 0, то не запускается
	CleanupInterval time.Duration
	// LimitEntries лимит записей в кэше, при достижении которого начинается вытеснение, если 0, то неограничено
	LimitEntries int
	// DataStore файл для восстановления кэша при перезапуске приложения, если пусто, то не восстанавливает
	DataStore string
}
