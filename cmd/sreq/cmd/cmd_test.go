package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCmd(t *testing.T) {
	// Test that root command exists and has correct properties
	if rootCmd.Use != "sreq" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "sreq")
	}

	if rootCmd.Short == "" {
		t.Error("rootCmd.Short should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("rootCmd.Long should not be empty")
	}
}

func TestRootCmd_PersistentFlags(t *testing.T) {
	flags := []struct {
		name      string
		shorthand string
	}{
		{"service", "s"},
		{"context", "c"},
		{"env", "e"},
		{"region", "r"},
		{"project", "p"},
		{"app", "a"},
		{"verbose", "v"},
		{"dry-run", ""},
	}

	for _, f := range flags {
		t.Run(f.name, func(t *testing.T) {
			flag := rootCmd.PersistentFlags().Lookup(f.name)
			if flag == nil {
				t.Errorf("Flag --%s not found", f.name)
				return
			}
			if f.shorthand != "" && flag.Shorthand != f.shorthand {
				t.Errorf("Flag shorthand = %q, want %q", flag.Shorthand, f.shorthand)
			}
		})
	}
}

func TestVersionCmd(t *testing.T) {
	if versionCmd.Use != "version" {
		t.Errorf("versionCmd.Use = %q, want %q", versionCmd.Use, "version")
	}

	if versionCmd.Short == "" {
		t.Error("versionCmd should have short description")
	}
}

func TestVersionCmd_Run(t *testing.T) {
	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the version command
	versionCmd.Run(versionCmd, []string{})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Verify output contains version info
	if !strings.Contains(output, "sreq version") {
		t.Errorf("Version output should contain 'sreq version', got: %s", output)
	}
	if !strings.Contains(output, "go:") {
		t.Errorf("Version output should contain 'go:', got: %s", output)
	}
	if !strings.Contains(output, "os/arch:") {
		t.Errorf("Version output should contain 'os/arch:', got: %s", output)
	}
}

func TestSubCommands(t *testing.T) {
	// Get all subcommands
	subcommands := rootCmd.Commands()

	// Check expected commands exist
	expectedCmds := []string{"version", "init", "run", "service", "auth", "config", "env", "history", "cache", "tui"}

	cmdMap := make(map[string]bool)
	for _, cmd := range subcommands {
		cmdMap[cmd.Use] = true
		// Also check commands with arguments like "run [method] [path]"
		parts := strings.Fields(cmd.Use)
		if len(parts) > 0 {
			cmdMap[parts[0]] = true
		}
	}

	for _, expected := range expectedCmds {
		if !cmdMap[expected] {
			t.Errorf("Expected subcommand %q not found", expected)
		}
	}
}

func TestInitCmd(t *testing.T) {
	cmd := findSubCommand("init")
	if cmd == nil {
		t.Fatal("init command not found")
	}

	if cmd.Short == "" {
		t.Error("init command should have short description")
	}
}

func TestRunCmd(t *testing.T) {
	cmd := findSubCommand("run")
	if cmd == nil {
		t.Fatal("run command not found")
	}

	// Check run command flags
	flags := []string{"header", "data", "timeout", "output"}
	for _, f := range flags {
		if flag := cmd.Flags().Lookup(f); flag == nil {
			t.Errorf("run command missing flag: %s", f)
		}
	}
}

func TestServiceCmd(t *testing.T) {
	cmd := findSubCommand("service")
	if cmd == nil {
		t.Fatal("service command not found")
	}

	// Check subcommands
	subCmds := cmd.Commands()
	subCmdMap := make(map[string]bool)
	for _, sub := range subCmds {
		parts := strings.Fields(sub.Use)
		if len(parts) > 0 {
			subCmdMap[parts[0]] = true
		}
	}

	expectedSubCmds := []string{"add", "list", "remove"}
	for _, expected := range expectedSubCmds {
		if !subCmdMap[expected] {
			t.Errorf("service subcommand %q not found", expected)
		}
	}
}

func TestAuthCmd(t *testing.T) {
	cmd := findSubCommand("auth")
	if cmd == nil {
		t.Fatal("auth command not found")
	}

	// Check subcommands
	subCmds := cmd.Commands()
	subCmdMap := make(map[string]bool)
	for _, sub := range subCmds {
		parts := strings.Fields(sub.Use)
		if len(parts) > 0 {
			subCmdMap[parts[0]] = true
		}
	}

	expectedSubCmds := []string{"consul", "aws"}
	for _, expected := range expectedSubCmds {
		if !subCmdMap[expected] {
			t.Errorf("auth subcommand %q not found", expected)
		}
	}
}

func TestConfigCmd(t *testing.T) {
	cmd := findSubCommand("config")
	if cmd == nil {
		t.Fatal("config command not found")
	}

	// Check subcommands
	subCmds := cmd.Commands()
	subCmdMap := make(map[string]bool)
	for _, sub := range subCmds {
		parts := strings.Fields(sub.Use)
		if len(parts) > 0 {
			subCmdMap[parts[0]] = true
		}
	}

	expectedSubCmds := []string{"show", "path", "test"}
	for _, expected := range expectedSubCmds {
		if !subCmdMap[expected] {
			t.Errorf("config subcommand %q not found", expected)
		}
	}
}

func TestHistoryCmd(t *testing.T) {
	cmd := findSubCommand("history")
	if cmd == nil {
		t.Fatal("history command not found")
	}

	// History is a single command with flags (not subcommands)
	// Check key flags
	flags := []string{"service", "env", "method", "all", "clear", "before", "curl", "httpie", "replay"}
	for _, f := range flags {
		if flag := cmd.Flags().Lookup(f); flag == nil {
			t.Errorf("history command missing flag: %s", f)
		}
	}
}

func TestCacheCmd(t *testing.T) {
	cmd := findSubCommand("cache")
	if cmd == nil {
		t.Fatal("cache command not found")
	}

	// Check subcommands
	subCmds := cmd.Commands()
	subCmdMap := make(map[string]bool)
	for _, sub := range subCmds {
		parts := strings.Fields(sub.Use)
		if len(parts) > 0 {
			subCmdMap[parts[0]] = true
		}
	}

	expectedSubCmds := []string{"clear", "status"}
	for _, expected := range expectedSubCmds {
		if !subCmdMap[expected] {
			t.Errorf("cache subcommand %q not found", expected)
		}
	}
}

func TestEnvCmd(t *testing.T) {
	cmd := findSubCommand("env")
	if cmd == nil {
		t.Fatal("env command not found")
	}

	// Check subcommands
	subCmds := cmd.Commands()
	subCmdMap := make(map[string]bool)
	for _, sub := range subCmds {
		parts := strings.Fields(sub.Use)
		if len(parts) > 0 {
			subCmdMap[parts[0]] = true
		}
	}

	expectedSubCmds := []string{"list", "switch", "current"}
	for _, expected := range expectedSubCmds {
		if !subCmdMap[expected] {
			t.Errorf("env subcommand %q not found", expected)
		}
	}
}

func TestTuiCmd(t *testing.T) {
	cmd := findSubCommand("tui")
	if cmd == nil {
		t.Fatal("tui command not found")
	}

	if cmd.Short == "" {
		t.Error("tui command should have short description")
	}
}

// Helper function to find a subcommand by name
func findSubCommand(name string) *cobra.Command {
	for _, cmd := range rootCmd.Commands() {
		parts := strings.Fields(cmd.Use)
		if len(parts) > 0 && parts[0] == name {
			return cmd
		}
	}
	return nil
}

func TestExecute(t *testing.T) {
	// Test that Execute function exists and doesn't panic
	// We can't fully test Execute without causing side effects
	// but we can verify it's callable
	_ = Execute // Just verify it compiles
}

func TestVersionVariable(t *testing.T) {
	// Version should have a default value
	if Version == "" {
		t.Error("Version should not be empty")
	}
}

// setupTestEnv creates a temporary HOME directory with test config
func setupTestEnv(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tmpDir)

	cleanup := func() {
		_ = os.Setenv("HOME", oldHome)
	}

	return tmpDir, cleanup
}

// createTestConfig creates a minimal config file for testing
func createTestConfig(t *testing.T, tmpDir string) {
	t.Helper()

	configDir := filepath.Join(tmpDir, ".sreq")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	configContent := `
providers: {}

environments:
  - dev
  - staging
  - prod

default_env: dev

services: {}
`
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	servicesContent := `
services: {}
`
	servicesPath := filepath.Join(configDir, "services.yaml")
	if err := os.WriteFile(servicesPath, []byte(servicesContent), 0644); err != nil {
		t.Fatalf("Failed to write services: %v", err)
	}
}

func TestRunInit_AlreadyInitialized(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	// Create existing config
	createTestConfig(t, tmpDir)

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run init command
	err := runInit(initCmd, []string{})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Errorf("runInit() error = %v", err)
	}

	if !strings.Contains(output, "already exists") {
		t.Errorf("Expected 'already exists' message, got: %s", output)
	}
}

func TestRunInit_Fresh(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run init command on fresh directory
	err := runInit(initCmd, []string{})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Errorf("runInit() error = %v", err)
	}

	if !strings.Contains(output, "Initializing") {
		t.Errorf("Expected 'Initializing' message, got: %s", output)
	}

	// Verify files were created
	configDir := filepath.Join(tmpDir, ".sreq")
	if _, err := os.Stat(filepath.Join(configDir, "config.yaml")); os.IsNotExist(err) {
		t.Error("config.yaml was not created")
	}
	if _, err := os.Stat(filepath.Join(configDir, "services.yaml")); os.IsNotExist(err) {
		t.Error("services.yaml was not created")
	}
	if _, err := os.Stat(filepath.Join(configDir, ".key")); os.IsNotExist(err) {
		t.Error(".key was not created")
	}
}

func TestRunConfigPath(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	createTestConfig(t, tmpDir)

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runConfigPath(configPathCmd, []string{})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Errorf("runConfigPath() error = %v", err)
	}

	if !strings.Contains(output, ".sreq") || !strings.Contains(output, "config.yaml") {
		t.Errorf("Expected config path output, got: %s", output)
	}
}

func TestRunConfigPath_WithEnvVar(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	customPath := filepath.Join(tmpDir, "custom-config.yaml")
	_ = os.WriteFile(customPath, []byte("providers: {}"), 0644)

	oldEnv := os.Getenv("SREQ_CONFIG")
	_ = os.Setenv("SREQ_CONFIG", customPath)
	defer func() {
		if oldEnv == "" {
			_ = os.Unsetenv("SREQ_CONFIG")
		} else {
			_ = os.Setenv("SREQ_CONFIG", oldEnv)
		}
	}()

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runConfigPath(configPathCmd, []string{})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Errorf("runConfigPath() error = %v", err)
	}

	if !strings.Contains(output, "SREQ_CONFIG") {
		t.Errorf("Expected SREQ_CONFIG mention, got: %s", output)
	}
}

func TestRunConfigShow(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	createTestConfig(t, tmpDir)

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runConfigShow(configShowCmd, []string{})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Errorf("runConfigShow() error = %v", err)
	}

	if !strings.Contains(output, "Configuration:") {
		t.Errorf("Expected 'Configuration:' in output, got: %s", output)
	}
	if !strings.Contains(output, "Providers:") {
		t.Errorf("Expected 'Providers:' in output, got: %s", output)
	}
	if !strings.Contains(output, "Environments:") {
		t.Errorf("Expected 'Environments:' in output, got: %s", output)
	}
}

func TestRunEnvList(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	createTestConfig(t, tmpDir)

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runEnvList(envListCmd, []string{})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Errorf("runEnvList() error = %v", err)
	}

	if !strings.Contains(output, "Available environments") {
		t.Errorf("Expected 'Available environments' in output, got: %s", output)
	}
	if !strings.Contains(output, "dev") {
		t.Errorf("Expected 'dev' in output, got: %s", output)
	}
	if !strings.Contains(output, "default") {
		t.Errorf("Expected 'default' marker in output, got: %s", output)
	}
}

func TestRunEnvCurrent(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	createTestConfig(t, tmpDir)

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runEnvCurrent(envCurrentCmd, []string{})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Errorf("runEnvCurrent() error = %v", err)
	}

	if !strings.Contains(output, "dev") {
		t.Errorf("Expected 'dev' in output, got: %s", output)
	}
}

func TestRunEnvSwitch(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	createTestConfig(t, tmpDir)

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runEnvSwitch(envSwitchCmd, []string{"staging"})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Errorf("runEnvSwitch() error = %v", err)
	}

	if !strings.Contains(output, "Switched") && !strings.Contains(output, "staging") {
		t.Errorf("Expected switch confirmation, got: %s", output)
	}
}

func TestRunEnvSwitch_InvalidEnv(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	createTestConfig(t, tmpDir)

	err := runEnvSwitch(envSwitchCmd, []string{"nonexistent"})

	if err == nil {
		t.Error("Expected error for invalid environment")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' in error, got: %v", err)
	}
}

func TestRunServiceList(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	createTestConfig(t, tmpDir)

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runServiceList(serviceListCmd, []string{})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Errorf("runServiceList() error = %v", err)
	}

	if !strings.Contains(output, "Configured services") {
		t.Errorf("Expected 'Configured services' in output, got: %s", output)
	}
}

func TestRunServiceAdd_SimpleMode(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	createTestConfig(t, tmpDir)

	// Set the flags
	consulKey = "auth"
	awsPrefix = ""
	pathMappings = nil

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runServiceAdd(serviceAddCmd, []string{"test-service"})

	_ = w.Close()
	os.Stdout = oldStdout

	// Reset flags
	consulKey = ""
	awsPrefix = ""
	pathMappings = nil

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Errorf("runServiceAdd() error = %v", err)
	}

	if !strings.Contains(output, "Added service") {
		t.Errorf("Expected 'Added service' in output, got: %s", output)
	}
	if !strings.Contains(output, "simple") {
		t.Errorf("Expected 'simple' mode in output, got: %s", output)
	}
}

func TestRunServiceAdd_AdvancedMode(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	createTestConfig(t, tmpDir)

	// Set the flags for advanced mode
	consulKey = ""
	awsPrefix = ""
	pathMappings = []string{"base_url=consul:services/test/url", "password=aws:secrets/test#password"}

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runServiceAdd(serviceAddCmd, []string{"advanced-service"})

	_ = w.Close()
	os.Stdout = oldStdout

	// Reset flags
	consulKey = ""
	awsPrefix = ""
	pathMappings = nil

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Errorf("runServiceAdd() error = %v", err)
	}

	if !strings.Contains(output, "Added service") {
		t.Errorf("Expected 'Added service' in output, got: %s", output)
	}
	if !strings.Contains(output, "advanced") {
		t.Errorf("Expected 'advanced' mode in output, got: %s", output)
	}
}

func TestRunServiceAdd_MixedMode(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	createTestConfig(t, tmpDir)

	// Set both simple and advanced flags (should fail)
	consulKey = "auth"
	pathMappings = []string{"base_url=test"}

	err := runServiceAdd(serviceAddCmd, []string{"mixed-service"})

	// Reset flags
	consulKey = ""
	pathMappings = nil

	if err == nil {
		t.Error("Expected error for mixed mode flags")
	}
}

func TestRunServiceAdd_NoFlags(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	createTestConfig(t, tmpDir)

	// No flags set
	consulKey = ""
	awsPrefix = ""
	pathMappings = nil

	err := runServiceAdd(serviceAddCmd, []string{"no-flags-service"})

	if err == nil {
		t.Error("Expected error when no flags provided")
	}
}

func TestRunServiceRemove(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	createTestConfig(t, tmpDir)

	// First add a service
	consulKey = "test"
	err := runServiceAdd(serviceAddCmd, []string{"removable-service"})
	consulKey = ""
	if err != nil {
		t.Fatalf("Failed to add service: %v", err)
	}

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = runServiceRemove(serviceRemoveCmd, []string{"removable-service"})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Errorf("runServiceRemove() error = %v", err)
	}

	if !strings.Contains(output, "Removed service") {
		t.Errorf("Expected 'Removed service' in output, got: %s", output)
	}
}

func TestRunServiceRemove_NotFound(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	createTestConfig(t, tmpDir)

	err := runServiceRemove(serviceRemoveCmd, []string{"nonexistent-service"})

	if err == nil {
		t.Error("Expected error for nonexistent service")
	}
}

func TestGenerateEncryptionKey(t *testing.T) {
	tmpDir := t.TempDir()

	// First call should create key
	err := generateEncryptionKey(tmpDir)
	if err != nil {
		t.Fatalf("generateEncryptionKey() error = %v", err)
	}

	keyPath := filepath.Join(tmpDir, ".key")
	info, err := os.Stat(keyPath)
	if os.IsNotExist(err) {
		t.Fatal(".key file was not created")
	}

	if info.Size() != 32 {
		t.Errorf("Key size = %d, want 32", info.Size())
	}

	// Check permissions (should be 0600)
	if info.Mode().Perm() != 0600 {
		t.Errorf("Key permissions = %o, want 0600", info.Mode().Perm())
	}

	// Second call should not overwrite
	originalKey, _ := os.ReadFile(keyPath)
	err = generateEncryptionKey(tmpDir)
	if err != nil {
		t.Fatalf("generateEncryptionKey() second call error = %v", err)
	}

	newKey, _ := os.ReadFile(keyPath)
	if string(originalKey) != string(newKey) {
		t.Error("Existing key was overwritten")
	}
}

func TestLoadServicesConfig_NotFound(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	// Don't create config, just the directory
	configDir := filepath.Join(tmpDir, ".sreq")
	_ = os.MkdirAll(configDir, 0755)

	_, err := LoadServicesConfig()

	if err == nil {
		t.Error("Expected error for missing services file")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' in error, got: %v", err)
	}
}

func TestLoadServicesConfig_Valid(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	createTestConfig(t, tmpDir)

	// Add a service to test
	configDir := filepath.Join(tmpDir, ".sreq")
	servicesContent := `
services:
  test-service:
    consul_key: test
`
	_ = os.WriteFile(filepath.Join(configDir, "services.yaml"), []byte(servicesContent), 0644)

	services, err := LoadServicesConfig()

	if err != nil {
		t.Fatalf("LoadServicesConfig() error = %v", err)
	}

	if services == nil {
		t.Fatal("Expected services map, got nil")
	}
}

func TestRunCacheStatus_CacheDisabled(t *testing.T) {
	// Set CI environment to disable cache
	oldCI := os.Getenv("CI")
	_ = os.Setenv("CI", "true")
	defer func() {
		if oldCI == "" {
			_ = os.Unsetenv("CI")
		} else {
			_ = os.Setenv("CI", oldCI)
		}
	}()

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runCacheStatus(cacheStatusCmd, []string{})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Errorf("runCacheStatus() error = %v", err)
	}

	if !strings.Contains(output, "disabled") {
		t.Errorf("Expected 'disabled' in output, got: %s", output)
	}
}

func TestRunCacheStatus_NotInitialized(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	// Create config but no .key file
	configDir := filepath.Join(tmpDir, ".sreq")
	_ = os.MkdirAll(configDir, 0755)
	_ = os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte("providers: {}"), 0644)

	// Unset CI to enable cache
	oldCI := os.Getenv("CI")
	_ = os.Unsetenv("CI")
	oldNoCache := os.Getenv("SREQ_NO_CACHE")
	_ = os.Unsetenv("SREQ_NO_CACHE")
	defer func() {
		if oldCI != "" {
			_ = os.Setenv("CI", oldCI)
		}
		if oldNoCache != "" {
			_ = os.Setenv("SREQ_NO_CACHE", oldNoCache)
		}
	}()

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runCacheStatus(cacheStatusCmd, []string{})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Errorf("runCacheStatus() error = %v", err)
	}

	if !strings.Contains(output, "not initialized") {
		t.Errorf("Expected 'not initialized' in output, got: %s", output)
	}
}

func TestRunCacheClear_NotInitialized(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	// Create config but no .key file
	configDir := filepath.Join(tmpDir, ".sreq")
	_ = os.MkdirAll(configDir, 0755)
	_ = os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte("providers: {}"), 0644)

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runCacheClear(cacheClearCmd, []string{})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Errorf("runCacheClear() error = %v", err)
	}

	if !strings.Contains(output, "not initialized") {
		t.Errorf("Expected 'not initialized' in output, got: %s", output)
	}
}

func TestSyncCmd_CacheDisabled(t *testing.T) {
	// Set CI environment to disable cache
	oldCI := os.Getenv("CI")
	_ = os.Setenv("CI", "true")
	defer func() {
		if oldCI == "" {
			_ = os.Unsetenv("CI")
		} else {
			_ = os.Setenv("CI", oldCI)
		}
	}()

	err := runSync(syncCmd, []string{"dev"})

	if err == nil {
		t.Error("Expected error when cache is disabled")
	}

	if !strings.Contains(err.Error(), "disabled") {
		t.Errorf("Expected 'disabled' in error, got: %v", err)
	}
}

func TestServiceAddCmd_Flags(t *testing.T) {
	// Test that service add command has correct flags
	flags := []string{"consul-key", "aws-prefix", "path"}
	for _, f := range flags {
		if flag := serviceAddCmd.Flags().Lookup(f); flag == nil {
			t.Errorf("service add command missing flag: %s", f)
		}
	}
}

func TestSyncCmd_Flags(t *testing.T) {
	// Test that sync command has correct flags
	flags := []string{"all", "force"}
	for _, f := range flags {
		if flag := syncCmd.Flags().Lookup(f); flag == nil {
			t.Errorf("sync command missing flag: %s", f)
		}
	}
}
