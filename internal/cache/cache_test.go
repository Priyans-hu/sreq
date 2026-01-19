package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Priyans-hu/sreq/pkg/types"
)

func setupTestCache(t *testing.T) (*Cache, string) {
	tmpDir, err := os.MkdirTemp("", "sreq-cache-test")
	if err != nil {
		t.Fatal(err)
	}

	// Generate and save key
	key, err := GenerateKey()
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		t.Fatal(err)
	}
	if err := SaveKey(tmpDir, key); err != nil {
		_ = os.RemoveAll(tmpDir)
		t.Fatal(err)
	}

	c, err := New(Config{
		ConfigDir: tmpDir,
		TTL:       1 * time.Hour,
	})
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		t.Fatal(err)
	}

	return c, tmpDir
}

func TestCache_SetAndGet(t *testing.T) {
	c, tmpDir := setupTestCache(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	creds := &types.ResolvedCredentials{
		BaseURL:  "https://api.example.com",
		Username: "user",
		Password: "secret",
		APIKey:   "key123",
	}

	// Set
	if err := c.Set("auth-service", "dev", creds); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get
	got, err := c.Get("auth-service", "dev")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got == nil {
		t.Fatal("Get returned nil")
	}

	if got.BaseURL != creds.BaseURL {
		t.Errorf("BaseURL = %q, want %q", got.BaseURL, creds.BaseURL)
	}
	if got.Username != creds.Username {
		t.Errorf("Username = %q, want %q", got.Username, creds.Username)
	}
	if got.Password != creds.Password {
		t.Errorf("Password = %q, want %q", got.Password, creds.Password)
	}
}

func TestCache_CacheMiss(t *testing.T) {
	c, tmpDir := setupTestCache(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Get non-existent entry
	got, err := c.Get("nonexistent", "dev")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got != nil {
		t.Error("Expected nil for cache miss")
	}
}

func TestCache_Expiry(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sreq-cache-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Generate and save key
	key, _ := GenerateKey()
	_ = SaveKey(tmpDir, key)

	// Create cache with very short TTL
	c, err := New(Config{
		ConfigDir: tmpDir,
		TTL:       1 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}

	creds := &types.ResolvedCredentials{
		BaseURL: "https://api.example.com",
	}

	// Set
	if err := c.Set("auth-service", "dev", creds); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Wait for expiry
	time.Sleep(10 * time.Millisecond)

	// Get should return nil (expired)
	got, err := c.Get("auth-service", "dev")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got != nil {
		t.Error("Expected nil for expired entry")
	}
}

func TestCache_Delete(t *testing.T) {
	c, tmpDir := setupTestCache(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	creds := &types.ResolvedCredentials{
		BaseURL: "https://api.example.com",
	}

	// Set
	_ = c.Set("auth-service", "dev", creds)

	// Delete
	if err := c.Delete("auth-service", "dev"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Get should return nil
	got, _ := c.Get("auth-service", "dev")
	if got != nil {
		t.Error("Expected nil after delete")
	}
}

func TestCache_Clear(t *testing.T) {
	c, tmpDir := setupTestCache(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	creds := &types.ResolvedCredentials{
		BaseURL: "https://api.example.com",
	}

	// Set multiple
	_ = c.Set("auth-service", "dev", creds)
	_ = c.Set("billing-service", "dev", creds)
	_ = c.Set("auth-service", "prod", creds)

	// Clear all
	if err := c.Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// All should return nil
	for _, svc := range []string{"auth-service", "billing-service"} {
		for _, env := range []string{"dev", "prod"} {
			got, _ := c.Get(svc, env)
			if got != nil {
				t.Errorf("Expected nil for %s/%s after clear", svc, env)
			}
		}
	}
}

func TestCache_ClearEnv(t *testing.T) {
	c, tmpDir := setupTestCache(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	creds := &types.ResolvedCredentials{
		BaseURL: "https://api.example.com",
	}

	// Set for different envs
	_ = c.Set("auth-service", "dev", creds)
	_ = c.Set("auth-service", "prod", creds)

	// Clear only dev
	if err := c.ClearEnv("dev"); err != nil {
		t.Fatalf("ClearEnv failed: %v", err)
	}

	// Dev should be gone
	got, _ := c.Get("auth-service", "dev")
	if got != nil {
		t.Error("Expected nil for dev after ClearEnv")
	}

	// Prod should still exist
	got, _ = c.Get("auth-service", "prod")
	if got == nil {
		t.Error("Expected prod to still exist")
	}
}

func TestCache_Status(t *testing.T) {
	c, tmpDir := setupTestCache(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	creds := &types.ResolvedCredentials{
		BaseURL: "https://api.example.com",
	}

	_ = c.Set("auth-service", "dev", creds)
	_ = c.Set("billing-service", "dev", creds)

	status, err := c.Status()
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}

	if status.EntryCount != 2 {
		t.Errorf("EntryCount = %d, want 2", status.EntryCount)
	}
	if len(status.Entries) != 2 {
		t.Errorf("len(Entries) = %d, want 2", len(status.Entries))
	}
}

func TestEncryption_RoundTrip(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	plaintext := []byte("secret data that needs encryption")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Ciphertext should be different from plaintext
	if string(ciphertext) == string(plaintext) {
		t.Error("Ciphertext equals plaintext")
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("Decrypted = %q, want %q", decrypted, plaintext)
	}
}

func TestEncryption_WrongKey(t *testing.T) {
	key1, _ := GenerateKey()
	key2, _ := GenerateKey()

	plaintext := []byte("secret data")
	ciphertext, _ := Encrypt(key1, plaintext)

	// Decrypt with wrong key should fail
	_, err := Decrypt(key2, ciphertext)
	if err == nil {
		t.Error("Expected error when decrypting with wrong key")
	}
}

func TestKeyExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sreq-cache-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Should not exist initially
	if KeyExists(tmpDir) {
		t.Error("Key should not exist initially")
	}

	// Generate and save key
	key, _ := GenerateKey()
	_ = SaveKey(tmpDir, key)

	// Should exist now
	if !KeyExists(tmpDir) {
		t.Error("Key should exist after save")
	}
}

func TestLoadKey_NotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sreq-cache-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	_, err = LoadKey(tmpDir)
	if err == nil {
		t.Error("Expected error when key doesn't exist")
	}
}

func TestIsDisabled(t *testing.T) {
	// Save original values and restore after test
	origNoCache := os.Getenv("SREQ_NO_CACHE")
	origCI := os.Getenv("CI")
	defer func() {
		_ = os.Setenv("SREQ_NO_CACHE", origNoCache)
		_ = os.Setenv("CI", origCI)
	}()

	// Clear env vars
	_ = os.Unsetenv("SREQ_NO_CACHE")
	_ = os.Unsetenv("CI")

	if IsDisabled() {
		t.Error("Should not be disabled when env vars are not set")
	}

	// Test SREQ_NO_CACHE
	_ = os.Setenv("SREQ_NO_CACHE", "1")
	if !IsDisabled() {
		t.Error("Should be disabled when SREQ_NO_CACHE=1")
	}
	_ = os.Unsetenv("SREQ_NO_CACHE")

	// Test CI=true
	_ = os.Setenv("CI", "true")
	if !IsDisabled() {
		t.Error("Should be disabled when CI=true")
	}
}

func TestEntry_IsExpired(t *testing.T) {
	// Non-expired
	e := Entry{
		CachedAt:   time.Now(),
		TTLSeconds: 3600,
	}
	if e.IsExpired() {
		t.Error("Should not be expired")
	}

	// Expired
	e = Entry{
		CachedAt:   time.Now().Add(-2 * time.Hour),
		TTLSeconds: 3600,
	}
	if !e.IsExpired() {
		t.Error("Should be expired")
	}
}

func TestCacheFilePermissions(t *testing.T) {
	c, tmpDir := setupTestCache(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	creds := &types.ResolvedCredentials{
		BaseURL: "https://api.example.com",
	}

	_ = c.Set("auth-service", "dev", creds)

	// Check file permissions
	cacheFile := filepath.Join(tmpDir, "cache", "dev", "auth-service-dev.enc")
	info, err := os.Stat(cacheFile)
	if err != nil {
		t.Fatalf("Failed to stat cache file: %v", err)
	}

	mode := info.Mode().Perm()
	if mode != 0600 {
		t.Errorf("Cache file permissions = %o, want 0600", mode)
	}
}
