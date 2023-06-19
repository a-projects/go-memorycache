package memorycache

import (
	"context"
	"encoding/gob"
	"math"
	"os"
	"sync"
	"time"
)

// MemoryCache кэш
type MemoryCache struct {
	options MemoryCacheOptions
	store   map[string]memoryCacheEntry
	mutex   sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
}

// Get предоставляет значение записи по ключу
func (m *MemoryCache) Get(key string) (value interface{}, ok bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	entry, ok := m.store[key]

	if ok && time.Since(entry.Expiration) >= 0 {
		return nil, false
	}

	return entry.Value, ok
}

// Set добавляет или изменяет значение записи по ключу
func (m *MemoryCache) Set(key string, value interface{}, options MemoryCacheEntryOptions) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// вытеснение слабой записи из кэша при достижении лимита
	if m.options.LimitEntries != 0 && len(m.store) >= m.options.LimitEntries {
		if _, ok := m.store[key]; !ok {
			now := time.Now()
			var wimpEntry string
			var wimpEntryDurability MemoryCacheEntryDurability = math.MaxInt

			for key, entry := range m.store {
				if now.Sub(entry.Expiration) >= 0 {
					wimpEntry = key
					break
				}

				if entry.Durability < wimpEntryDurability {
					wimpEntryDurability = entry.Durability
					wimpEntry = key
				}
			}

			delete(m.store, wimpEntry)
		}
	}

	m.store[key] = memoryCacheEntry{
		Value:      value,
		Expiration: options.Expiration,
		Durability: options.Durability,
	}
}

// Del удаляет запись по ключу
func (m *MemoryCache) Del(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.store, key)
}

// Count предоставляет количество записей в кэше
func (m *MemoryCache) Count() (count int) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.store)
}

// Reset удаляет все записи из кэша
func (m *MemoryCache) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.store = make(map[string]memoryCacheEntry)
}

// startCleanup запускает механизм очистки
func (m *MemoryCache) startCleanup() {
	for {
		select {
		case <-m.ctx.Done(): // ожидаем остановку сервиса
			return
		case <-time.After(m.options.CleanupInterval): // ожидаем таймаут очистки
			m.сleanup()
		}
	}
}

// сleanup очищает от записей с истёкшим временем жизни
func (m *MemoryCache) сleanup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for key, entry := range m.store {
		if time.Since(entry.Expiration) >= 0 {
			delete(m.store, key)
		}
	}
}

// New конструктор
// также если задан DataStore, то загружает данные кэша из файла
func New(options MemoryCacheOptions) (memorycache *MemoryCache) {
	ctx, cancel := context.WithCancel(context.Background())

	memorycache = &MemoryCache{
		store:   make(map[string]memoryCacheEntry),
		options: options,
		ctx:     ctx,
		cancel:  cancel,
	}

	if options.DataStore != "" {
		memorycache.load()
	}

	if options.CleanupInterval > 0 {
		go memorycache.startCleanup()
	}

	return memorycache
}

// Close деструктор
// также если задан DataStore, то сохраняет данные кэша в файл
func (m *MemoryCache) Close() {
	if m.cancel != nil {
		m.cancel()
	}

	if m.options.DataStore != "" {
		m.save()
	}
}

// load загружает данные кэша из файла
func (m *MemoryCache) load() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	file, err := os.Open(m.options.DataStore)

	if err == nil {
		defer file.Close()
		decoder := gob.NewDecoder(file)
		decoder.Decode(&m.store)
	}
}

// save сохраняет данные кэша в файл
func (m *MemoryCache) save() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	file, err := os.Create(m.options.DataStore)

	if err == nil {
		defer file.Close()
		encoder := gob.NewEncoder(file)
		encoder.Encode(m.store)
	}
}
