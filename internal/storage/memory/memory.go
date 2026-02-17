package memory

import (
	"sync"
	"urlshort/internal/generator"
	"urlshort/internal/storage"
)

type MemoryStorage struct{
	mu sync.RWMutex
	shortToLong map[string]string
	longToShort map[string]string
}

func NewMemoryStorage() *MemoryStorage{
	return &MemoryStorage{
		shortToLong: make(map[string]string),
		longToShort: make(map[string]string),
	}
}

func (mem *MemoryStorage) Save(origURL string) (string, error) {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	if short, ok := mem.longToShort[origURL]; ok{
		return short, nil
	}

	for {
		short := generator.Generate()
		if _, exist := mem.shortToLong[short]; !exist{
			mem.shortToLong[short] = origURL
			mem.longToShort[origURL] = short
			return short, nil
		}
	}
}

func (mem *MemoryStorage) Get(shortURL string) (string, error) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	if original, ok := mem.shortToLong[shortURL]; ok {
		return original, nil
	}
	return "", storage.ErrNotFound
}

func (mem *MemoryStorage) Close() error {
	return nil
}