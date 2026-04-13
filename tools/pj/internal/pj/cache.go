package pj

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func loadCache(path string) (*Cache, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cache Cache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("decode cache: %w", err)
	}

	return &cache, nil
}

func loadCacheRequired(path string) (*Cache, error) {
	cache, err := loadCache(path)
	if err == nil {
		return cache, nil
	}
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("cache not found at %s; run `pj sync` first", path)
	}
	return nil, err
}

func writeCache(path string, cache *Cache) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("encode cache: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write cache: %w", err)
	}
	return nil
}
