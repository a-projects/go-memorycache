package memorycache

import (
	"context"
	"encoding/gob"
	"math"
	"os"
	"sync"
	"time"
)

// MemoryCache cache
type MemoryCache struct {
	options MemoryCacheOptions
	store   map[string]memoryCacheEntry
	mutex   sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
}

// Get provide entry value by key
func (m *MemoryCache) Get(key string) (value interface{}, ok bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	entry, ok := m.store[key]

	if ok && time.Since(entry.Expiration) >= 0 {
		return nil, false
	}

	return entry.Value, ok
}

// Set add entry or update value by key
func (m *MemoryCache) Set(key string, value interface{}, options MemoryCacheEntryOptions) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// удаление слабой записи из кеша при достижении лимита
	if m.options.LimitEntries != 0 && len(m.store) >= m.options.LimitEntries {
		if _, ok := m.store[key]; !ok {
			now := time.Now()
			var wimpEntry string
			var wimpEntryExpiration time.Time = time.Now()
			var wimpEntryDurability MemoryCacheEntryDurability = math.MaxInt

			// алгоритм, исходит из того, что всегда получит запись т.к. сравнивает стойкость с MaxInt
			for key, entry := range m.store {
				if now.Sub(entry.Expiration) >= 0 {
					wimpEntry = key
					break
				}

				if entry.Durability < wimpEntryDurability {
					wimpEntryDurability = entry.Durability
					wimpEntryExpiration = entry.Expiration
					wimpEntry = key
					continue
				}

				if entry.Durability == wimpEntryDurability {
					if wimpEntryExpiration.Sub(entry.Expiration) > 0 {
						wimpEntryDurability = entry.Durability
						wimpEntryExpiration = entry.Expiration
						wimpEntry = key
					}
				}
			}

			delete(m.store, wimpEntry)
		}
	}

	if options.Lifetime > 0 {
		m.store[key] = memoryCacheEntry{
			Value:      value,
			Expiration: time.Now().Add(options.Lifetime),
			Durability: options.Durability,
		}
	}
}

// Del delete entry by key
func (m *MemoryCache) Del(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.store, key)
}

// Count provide count entries in cache
func (m *MemoryCache) Count() (count int) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.store)
}

// Reset remove all entries from cache
func (m *MemoryCache) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.store = make(map[string]memoryCacheEntry)
}

// startCleanup start cleaning mechanism
func (m *MemoryCache) startCleanup() {
	for {
		select {
		case <-m.ctx.Done(): // ожидать остановки сервиса
			return
		case <-time.After(m.options.CleanupInterval): // ожидать таймаут очистки
			m.сleanup() //гонки не будет т.к. cleanup в той же горутине, что и time.After
		}
	}
}

// сleanup cleaning up expired records
func (m *MemoryCache) сleanup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.store != nil {
		for key, entry := range m.store {
			if time.Since(entry.Expiration) >= 0 {
				delete(m.store, key)
			}
		}
	}
}

// New construct and load cache data from file if set FileName
func New(options MemoryCacheOptions) (memorycache *MemoryCache) {
	ctx, cancel := context.WithCancel(context.Background())

	memorycache = &MemoryCache{
		store:   make(map[string]memoryCacheEntry),
		options: options,
		ctx:     ctx,
		cancel:  cancel,
	}

	if options.FileName != "" {
		memorycache.load()
	}

	if options.CleanupInterval > 0 {
		go memorycache.startCleanup()
	}

	return memorycache
}

// Close destruct and save cache data in file if set FileName
func (m *MemoryCache) Close() {
	if m.cancel != nil {
		m.cancel()
	}

	if m.options.FileName != "" {
		m.save()
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.store = nil
}

// load load cache data from file
func (m *MemoryCache) load() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	file, err := os.Open(m.options.FileName)

	if err == nil {
		defer file.Close()
		decoder := gob.NewDecoder(file)
		decoder.Decode(&m.store)
	}
}

// save save cache data to file
func (m *MemoryCache) save() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	file, err := os.Create(m.options.FileName)

	if err == nil {
		defer file.Close()
		encoder := gob.NewEncoder(file)
		encoder.Encode(m.store)
	}
}
