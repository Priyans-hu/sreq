package dotenv

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Priyans-hu/sreq/internal/providers"
)

// Provider implements the providers.Provider interface for .env files
type Provider struct {
	files  []string
	paths  map[string]string
	values map[string]string
	mu     sync.RWMutex
	loaded bool
}

// Config holds dotenv provider configuration
type Config struct {
	// Files is a list of .env files to load (in order, later files override earlier)
	// Supports: ".env", ".env.local", ".env.{env}", etc.
	Files []string
	// File is a single file path (for backward compatibility)
	File string
	// Paths contains path templates for credential resolution
	Paths map[string]string
}

// New creates a new dotenv provider
func New(cfg Config) (*Provider, error) {
	files := cfg.Files
	if cfg.File != "" {
		files = append([]string{cfg.File}, files...)
	}

	// Default to .env if no files specified
	if len(files) == 0 {
		files = []string{".env"}
	}

	return &Provider{
		files:  files,
		paths:  cfg.Paths,
		values: make(map[string]string),
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "dotenv"
}

// loadFiles loads all configured .env files
func (p *Provider) loadFiles() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.loaded {
		return nil
	}

	for _, file := range p.files {
		// Expand home directory
		if strings.HasPrefix(file, "~/") {
			home, err := os.UserHomeDir()
			if err == nil {
				file = filepath.Join(home, file[2:])
			}
		}

		// Try to load the file
		if err := p.loadFile(file); err != nil {
			// Only error if it's not a "file not found" error
			if !os.IsNotExist(err) {
				return err
			}
			// File doesn't exist, skip silently
		}
	}

	p.loaded = true
	return nil
}

// loadFile parses a single .env file
func (p *Provider) loadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		key, value, err := parseLine(line)
		if err != nil {
			// Skip malformed lines silently
			continue
		}

		p.values[key] = value
	}

	return scanner.Err()
}

// parseLine parses a single line from a .env file
// Supports formats:
// - KEY=value
// - KEY="value with spaces"
// - KEY='value with spaces'
// - export KEY=value
func parseLine(line string) (string, string, error) {
	// Remove 'export ' prefix if present
	line = strings.TrimPrefix(line, "export ")
	line = strings.TrimSpace(line)

	// Find the first '='
	idx := strings.Index(line, "=")
	if idx == -1 {
		return "", "", fmt.Errorf("invalid format: no '=' found")
	}

	key := strings.TrimSpace(line[:idx])
	value := strings.TrimSpace(line[idx+1:])

	// Validate key
	if key == "" {
		return "", "", fmt.Errorf("empty key")
	}

	// Remove quotes from value
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}
	}

	// Handle escape sequences in double-quoted strings
	value = strings.ReplaceAll(value, "\\n", "\n")
	value = strings.ReplaceAll(value, "\\t", "\t")
	value = strings.ReplaceAll(value, "\\\"", "\"")

	return key, value, nil
}

// Get retrieves a value from the loaded .env files
func (p *Provider) Get(ctx context.Context, key string) (string, error) {
	// Ensure files are loaded
	if err := p.loadFiles(); err != nil {
		return "", fmt.Errorf("failed to load .env files: %w", err)
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	// Convert key to uppercase (env vars are typically uppercase)
	envKey := strings.ToUpper(key)

	value, exists := p.values[envKey]
	if !exists {
		// Try original case as fallback
		value, exists = p.values[key]
	}
	if !exists {
		return "", fmt.Errorf("key '%s' not found in .env files", key)
	}

	return value, nil
}

// GetMultiple retrieves multiple values from the .env files
func (p *Provider) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	results := make(map[string]string)

	for _, key := range keys {
		value, err := p.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		results[key] = value
	}

	return results, nil
}

// GetWithTemplate retrieves a value using a path template
func (p *Provider) GetWithTemplate(ctx context.Context, template string, vars map[string]string) (string, error) {
	key := ResolvePath(template, vars)
	return p.Get(ctx, key)
}

// Health checks if at least one .env file exists
func (p *Provider) Health(ctx context.Context) error {
	for _, file := range p.files {
		// Expand home directory
		if strings.HasPrefix(file, "~/") {
			home, err := os.UserHomeDir()
			if err == nil {
				file = filepath.Join(home, file[2:])
			}
		}

		if _, err := os.Stat(file); err == nil {
			return nil
		}
	}
	return fmt.Errorf("no .env files found: %v", p.files)
}

// Reload forces a reload of all .env files
func (p *Provider) Reload() error {
	p.mu.Lock()
	p.loaded = false
	p.values = make(map[string]string)
	p.mu.Unlock()

	return p.loadFiles()
}

// GetAll returns all loaded key-value pairs
func (p *Provider) GetAll() (map[string]string, error) {
	if err := p.loadFiles(); err != nil {
		return nil, err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make(map[string]string, len(p.values))
	for k, v := range p.values {
		result[k] = v
	}
	return result, nil
}

// ResolvePath replaces placeholders in a path template (case-insensitive)
func ResolvePath(template string, vars map[string]string) string {
	result := template
	for key, value := range vars {
		// Convert value to env-friendly format (uppercase, underscores)
		envValue := strings.ToUpper(value)
		envValue = strings.ReplaceAll(envValue, "-", "_")
		envValue = strings.ReplaceAll(envValue, ".", "_")
		// Replace both lowercase and uppercase placeholders
		result = strings.ReplaceAll(result, "{"+key+"}", envValue)
		result = strings.ReplaceAll(result, "{"+strings.ToUpper(key)+"}", envValue)
	}
	return result
}

// Ensure Provider implements the interface
var _ providers.Provider = (*Provider)(nil)
