//cache.go
package main

import (
    "context"
    "encoding/json"
    "errors"
    "os"
    "path/filepath"
    "sync"
    "time"

    "github.com/go-redis/redis/v8"
)

// CacheEntry stores the weather data and its expiration time (file cache)
type CacheEntry struct {
    Data      WeatherData `json:"data"`
    ExpiresAt time.Time   `json:"expires_at"`
}

// Cache defines the behaviour for caching weather data
type Cache interface {
    Get(cityKey string) (*WeatherData, error)
    Set(cityKey string, data *WeatherData) error
    Close() error
}

// -------------------------------------------------------------------
// FileCache – persistent JSON file cache
// -------------------------------------------------------------------
type FileCache struct {
    filePath string
    mu       sync.RWMutex
    store    map[string]CacheEntry
}

func NewFileCache() (*FileCache, error) {
    cacheDir, err := os.UserCacheDir()
    if err != nil {
        return nil, err
    }
    cacheDir = filepath.Join(cacheDir, "weather_cli")
    if err := os.MkdirAll(cacheDir, 0755); err != nil {
        return nil, err
    }
    filePath := filepath.Join(cacheDir, "cache.json")

    fc := &FileCache{
        filePath: filePath,
        store:    make(map[string]CacheEntry),
    }
    // load existing cache from disk
    data, err := os.ReadFile(filePath)
    if err == nil {
        _ = json.Unmarshal(data, &fc.store)
    }
    return fc, nil
}

func (fc *FileCache) Get(cityKey string) (*WeatherData, error) {
    fc.mu.RLock()
    defer fc.mu.RUnlock()
    entry, exists := fc.store[cityKey]
    if !exists || time.Now().After(entry.ExpiresAt) {
        // expired or not found
        return nil, errors.New("cache miss or expired")
    }
    return &entry.Data, nil
}

func (fc *FileCache) Set(cityKey string, data *WeatherData) error {
    fc.mu.Lock()
    defer fc.mu.Unlock()
    fc.store[cityKey] = CacheEntry{
        Data:      *data,
        ExpiresAt: time.Now().Add(cacheTTL),
    }
    // write to disk
    raw, err := json.MarshalIndent(fc.store, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(fc.filePath, raw, 0644)
}

func (fc *FileCache) Close() error { return nil }

// -------------------------------------------------------------------
// RedisCache – Redis backend (used when REDIS_ADDR env var is set)
// -------------------------------------------------------------------
type RedisCache struct {
    client *redis.Client
    ctx    context.Context
}

func NewRedisCache(addr string) (*RedisCache, error) {
    opt, err := redis.ParseURL(addr)
    if err != nil {
        // if not a URL, assume "host:port"
        opt = &redis.Options{Addr: addr}
    }
    client := redis.NewClient(opt)
    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, err
    }
    return &RedisCache{client: client, ctx: ctx}, nil
}

func (rc *RedisCache) Get(cityKey string) (*WeatherData, error) {
    val, err := rc.client.Get(rc.ctx, "weather:"+cityKey).Result()
    if err == redis.Nil {
        return nil, errors.New("cache miss")
    }
    if err != nil {
        return nil, err
    }
    var data WeatherData
    if err := json.Unmarshal([]byte(val), &data); err != nil {
        return nil, err
    }
    return &data, nil
}

func (rc *RedisCache) Set(cityKey string, data *WeatherData) error {
    raw, err := json.Marshal(data)
    if err != nil {
        return err
    }
    return rc.client.Set(rc.ctx, "weather:"+cityKey, raw, cacheTTL).Err()
}

func (rc *RedisCache) Close() error {
    return rc.client.Close()
}