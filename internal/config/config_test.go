package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
providers:
  consul:
    address: localhost:8500
    paths:
      base_url: "services/{service}/url"

environments:
  - dev
  - prod

default_env: dev

services:
  auth:
    consul_key: auth-service
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	// Verify providers
	if cfg.Providers["consul"].Address != "localhost:8500" {
		t.Errorf("Consul address = %q, want %q", cfg.Providers["consul"].Address, "localhost:8500")
	}

	// Verify environments
	if len(cfg.Environments) != 2 {
		t.Errorf("Environments count = %d, want 2", len(cfg.Environments))
	}

	// Verify default env
	if cfg.DefaultEnv != "dev" {
		t.Errorf("DefaultEnv = %q, want %q", cfg.DefaultEnv, "dev")
	}

	// Verify services
	if _, ok := cfg.Services["auth"]; !ok {
		t.Error("Service 'auth' not found")
	}
}

func TestLoadFromFile_NotFound(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/config.yaml")
	if err == nil {
		t.Error("LoadFromFile() expected error for nonexistent file")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Error should mention 'not found', got: %v", err)
	}
}

func TestLoadFromFile_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Write invalid YAML
	if err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	_, err := LoadFromFile(configPath)
	if err == nil {
		t.Error("LoadFromFile() expected error for invalid YAML")
	}
}

func TestLoadFromFile_WithServicesFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create main config
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
providers:
  consul:
    address: localhost:8500

services:
  existing:
    consul_key: existing-key
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Create services.yaml
	servicesPath := filepath.Join(tmpDir, "services.yaml")
	servicesContent := `
services:
  billing:
    consul_key: billing-key
  existing:
    consul_key: overridden-key
`
	if err := os.WriteFile(servicesPath, []byte(servicesContent), 0644); err != nil {
		t.Fatalf("failed to write services: %v", err)
	}

	cfg, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	// Check new service was added
	if _, ok := cfg.Services["billing"]; !ok {
		t.Error("Service 'billing' not found")
	}

	// Check existing service was overridden
	if cfg.Services["existing"].ConsulKey != "overridden-key" {
		t.Errorf("Service 'existing' consul_key = %q, want 'overridden-key'", cfg.Services["existing"].ConsulKey)
	}
}

func TestLoad_WithEnvVar(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "custom-config.yaml")

	configContent := `
providers:
  consul:
    address: custom:8500
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// Set SREQ_CONFIG env var
	oldEnv := os.Getenv("SREQ_CONFIG")
	_ = os.Setenv("SREQ_CONFIG", configPath)
	defer func() {
		if oldEnv == "" {
			_ = os.Unsetenv("SREQ_CONFIG")
		} else {
			_ = os.Setenv("SREQ_CONFIG", oldEnv)
		}
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Providers["consul"].Address != "custom:8500" {
		t.Errorf("Address = %q, want %q", cfg.Providers["consul"].Address, "custom:8500")
	}
}

func TestInit(t *testing.T) {
	// Use a temporary home directory
	tmpHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tmpHome)
	defer func() { _ = os.Setenv("HOME", oldHome) }()

	err := Init()
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Check config dir was created
	configDir := filepath.Join(tmpHome, DefaultConfigDir)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("Config directory was not created")
	}

	// Check config.yaml was created
	configPath := filepath.Join(configDir, DefaultConfigFile)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config.yaml was not created")
	}

	// Check services.yaml was created
	servicesPath := filepath.Join(configDir, DefaultServicesFile)
	if _, err := os.Stat(servicesPath); os.IsNotExist(err) {
		t.Error("services.yaml was not created")
	}

	// Verify config content is valid YAML
	cfg, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("Created config is invalid: %v", err)
	}
	if cfg.Providers == nil {
		t.Error("Config should have providers section")
	}
}

func TestInit_ExistingConfig(t *testing.T) {
	tmpHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tmpHome)
	defer func() { _ = os.Setenv("HOME", oldHome) }()

	// Create existing config
	configDir := filepath.Join(tmpHome, DefaultConfigDir)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	existingContent := "existing: content"
	configPath := filepath.Join(configDir, DefaultConfigFile)
	if err := os.WriteFile(configPath, []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to write existing config: %v", err)
	}

	// Init should not overwrite existing config
	err := Init()
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	if string(data) != existingContent {
		t.Error("Init() should not overwrite existing config")
	}
}

func TestGetConfigDir(t *testing.T) {
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir() error = %v", err)
	}

	if !strings.HasSuffix(dir, DefaultConfigDir) {
		t.Errorf("GetConfigDir() = %q, should end with %q", dir, DefaultConfigDir)
	}
}

func TestLoadFromFile_NilServices(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Config without services section
	configContent := `
providers:
  consul:
    address: localhost:8500
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	// Services map should be initialized
	if cfg.Services == nil {
		t.Error("Services should be initialized to empty map, not nil")
	}
}

func TestLoadFromFile_InvalidServicesYAML(t *testing.T) {
	tmpDir := t.TempDir()

	// Create main config
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("providers: {}"), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Create invalid services.yaml
	servicesPath := filepath.Join(tmpDir, "services.yaml")
	if err := os.WriteFile(servicesPath, []byte("invalid: yaml: ["), 0644); err != nil {
		t.Fatalf("failed to write services: %v", err)
	}

	_, err := LoadFromFile(configPath)
	if err == nil {
		t.Error("LoadFromFile() expected error for invalid services.yaml")
	}
}

func TestLoadFromFile_Contexts(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
providers:
  consul:
    address: localhost:8500

contexts:
  dev-us:
    project: myproject
    env: dev
    region: us-east-1
  prod-eu:
    project: myproject
    env: prod
    region: eu-west-1

default_context: dev-us
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	if len(cfg.Contexts) != 2 {
		t.Errorf("Contexts count = %d, want 2", len(cfg.Contexts))
	}

	if cfg.DefaultContext != "dev-us" {
		t.Errorf("DefaultContext = %q, want %q", cfg.DefaultContext, "dev-us")
	}

	devContext := cfg.Contexts["dev-us"]
	if devContext.Env != "dev" {
		t.Errorf("dev-us env = %q, want %q", devContext.Env, "dev")
	}
}
