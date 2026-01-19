package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Priyans-hu/sreq/pkg/types"
)

const (
	// DefaultCacheDir is the default cache directory name
	DefaultCacheDir = "cache"

	// DefaultTTL is the default cache TTL
	DefaultTTL = 1 * time.Hour

	// CacheFileExtension is the extension for cache files
	CacheFileExtension = ".enc"
)

// Entry represents a cached credential entry
type Entry struct {
	Service     string                   `json:"service"`
	Env         string                   `json:"env"`
	CachedAt    time.Time                `json:"cached_at"`
	TTLSeconds  int                      `json:"ttl_seconds"`
	Credentials *types.ResolvedCredentials `json:"credentials"`
}

// IsExpired checks if the cache entry has expired
func (e *Entry) IsExpired() bool {
	ttl := time.Duration(e.TTLSeconds) * time.Second
	return time.Since(e.CachedAt) > ttl
}

// ExpiresAt returns when the entry expires
func (e *Entry) ExpiresAt() time.Time {
	return e.CachedAt.Add(time.Duration(e.TTLSeconds) * time.Second)
}

// Cache manages credential caching
type Cache struct {
	configDir string
	cacheDir  string
	key       []byte
	ttl       time.Duration
}

// Config holds cache configuration
type Config struct {
	ConfigDir string
	TTL       time.Duration
}

// New creates a new cache manager
func New(cfg Config) (*Cache, error) {
	cacheDir := filepath.Join(cfg.ConfigDir, DefaultCacheDir)

	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Load encryption key
	key, err := LoadKey(cfg.ConfigDir)
	if err != nil {
		return nil, err
	}

	ttl := cfg.TTL
	if ttl == 0 {
		ttl = DefaultTTL
	}

	return &Cache{
		configDir: cfg.ConfigDir,
		cacheDir:  cacheDir,
		key:       key,
		ttl:       ttl,
	}, nil
}

// cacheFilePath returns the path to a cache file for a service/env
func (c *Cache) cacheFilePath(service, env string) string {
	filename := fmt.Sprintf("%s-%s%s", service, env, CacheFileExtension)
	return filepath.Join(c.cacheDir, env, filename)
}

// Get retrieves cached credentials for a service/env
func (c *Cache) Get(service, env string) (*types.ResolvedCredentials, error) {
	path := c.cacheFilePath(service, env)

	// Read encrypted file
	ciphertext, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Cache miss, not an error
		}
		return nil, fmt.Errorf("failed to read cache: %w", err)
	}

	// Decrypt
	plaintext, err := Decrypt(c.key, ciphertext)
	if err != nil {
		// Corrupted cache, remove it
		_ = os.Remove(path)
		return nil, nil
	}

	// Unmarshal
	var entry Entry
	if err := json.Unmarshal(plaintext, &entry); err != nil {
		// Corrupted cache, remove it
		_ = os.Remove(path)
		return nil, nil
	}

	// Check expiry
	if entry.IsExpired() {
		_ = os.Remove(path)
		return nil, nil
	}

	return entry.Credentials, nil
}

// Set caches credentials for a service/env
func (c *Cache) Set(service, env string, creds *types.ResolvedCredentials) error {
	entry := Entry{
		Service:     service,
		Env:         env,
		CachedAt:    time.Now(),
		TTLSeconds:  int(c.ttl.Seconds()),
		Credentials: creds,
	}

	// Marshal
	plaintext, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	// Encrypt
	ciphertext, err := Encrypt(c.key, plaintext)
	if err != nil {
		return fmt.Errorf("failed to encrypt cache: %w", err)
	}

	// Ensure env directory exists
	path := c.cacheFilePath(service, env)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Write with restricted permissions
	if err := os.WriteFile(path, ciphertext, 0600); err != nil {
		return fmt.Errorf("failed to write cache: %w", err)
	}

	return nil
}

// Delete removes cached credentials for a service/env
func (c *Cache) Delete(service, env string) error {
	path := c.cacheFilePath(service, env)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete cache: %w", err)
	}
	return nil
}

// Clear removes all cached credentials
func (c *Cache) Clear() error {
	if err := os.RemoveAll(c.cacheDir); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}
	// Recreate empty cache directory
	return os.MkdirAll(c.cacheDir, 0700)
}

// ClearEnv removes all cached credentials for an environment
func (c *Cache) ClearEnv(env string) error {
	envDir := filepath.Join(c.cacheDir, env)
	if err := os.RemoveAll(envDir); err != nil {
		return fmt.Errorf("failed to clear cache for env %s: %w", env, err)
	}
	return nil
}

// Status returns cache status information
type Status struct {
	Enabled    bool
	CacheDir   string
	TTL        time.Duration
	EntryCount int
	TotalSize  int64
	Entries    []EntryInfo
}

// EntryInfo contains information about a cache entry
type EntryInfo struct {
	Service   string
	Env       string
	CachedAt  time.Time
	ExpiresAt time.Time
	Size      int64
	Expired   bool
}

// Status returns cache status
func (c *Cache) Status() (*Status, error) {
	status := &Status{
		Enabled:  true,
		CacheDir: c.cacheDir,
		TTL:      c.ttl,
	}

	// Walk cache directory
	err := filepath.Walk(c.cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != CacheFileExtension {
			return nil
		}

		// Read and decrypt to get metadata
		ciphertext, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		plaintext, err := Decrypt(c.key, ciphertext)
		if err != nil {
			return nil
		}

		var entry Entry
		if err := json.Unmarshal(plaintext, &entry); err != nil {
			return nil
		}

		status.EntryCount++
		status.TotalSize += info.Size()
		status.Entries = append(status.Entries, EntryInfo{
			Service:   entry.Service,
			Env:       entry.Env,
			CachedAt:  entry.CachedAt,
			ExpiresAt: entry.ExpiresAt(),
			Size:      info.Size(),
			Expired:   entry.IsExpired(),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return status, nil
}

// IsDisabled checks if caching is disabled via environment variable
func IsDisabled() bool {
	// Check SREQ_NO_CACHE or CI environment
	if os.Getenv("SREQ_NO_CACHE") == "1" {
		return true
	}
	if os.Getenv("CI") == "true" || os.Getenv("CI") == "1" {
		return true
	}
	return false
}
