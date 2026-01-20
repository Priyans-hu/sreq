package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestSreqError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *SreqError
		contains []string
	}{
		{
			name: "message only",
			err: &SreqError{
				Message: "something went wrong",
			},
			contains: []string{"something went wrong"},
		},
		{
			name: "message with cause",
			err: &SreqError{
				Message: "failed to connect",
				Cause:   errors.New("connection refused"),
			},
			contains: []string{"failed to connect", "Cause:", "connection refused"},
		},
		{
			name: "message with suggestion",
			err: &SreqError{
				Message:    "config not found",
				Suggestion: "run init command",
			},
			contains: []string{"config not found", "Suggestion:", "run init command"},
		},
		{
			name: "full error",
			err: &SreqError{
				Type:       ErrConfig,
				Message:    "parse error",
				Cause:      errors.New("invalid yaml"),
				Suggestion: "check syntax",
			},
			contains: []string{"parse error", "Cause:", "invalid yaml", "Suggestion:", "check syntax"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := tt.err.Error()
			for _, s := range tt.contains {
				if !strings.Contains(errStr, s) {
					t.Errorf("Error() = %q, should contain %q", errStr, s)
				}
			}
		})
	}
}

func TestSreqError_Unwrap(t *testing.T) {
	cause := errors.New("original error")
	err := &SreqError{
		Message: "wrapped error",
		Cause:   cause,
	}

	if err.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
	}

	// Test nil cause
	errNoCause := &SreqError{Message: "no cause"}
	if errNoCause.Unwrap() != nil {
		t.Errorf("Unwrap() = %v, want nil", errNoCause.Unwrap())
	}
}

func TestConfigNotFound(t *testing.T) {
	err := ConfigNotFound("/path/to/config")
	if err.Type != ErrConfig {
		t.Errorf("Type = %v, want %v", err.Type, ErrConfig)
	}
	if !strings.Contains(err.Message, "/path/to/config") {
		t.Errorf("Message should contain path")
	}
	if err.Suggestion == "" {
		t.Error("Suggestion should not be empty")
	}
}

func TestConfigParseError(t *testing.T) {
	cause := errors.New("yaml error")
	err := ConfigParseError("/path/to/config", cause)
	if err.Type != ErrConfig {
		t.Errorf("Type = %v, want %v", err.Type, ErrConfig)
	}
	if err.Cause != cause {
		t.Errorf("Cause = %v, want %v", err.Cause, cause)
	}
}

func TestServiceNotFound(t *testing.T) {
	err := ServiceNotFound("my-service")
	if err.Type != ErrNotFound {
		t.Errorf("Type = %v, want %v", err.Type, ErrNotFound)
	}
	if !strings.Contains(err.Message, "my-service") {
		t.Errorf("Message should contain service name")
	}
	if !strings.Contains(err.Suggestion, "my-service") {
		t.Errorf("Suggestion should contain service name")
	}
}

func TestContextNotFound(t *testing.T) {
	err := ContextNotFound("prod-us")
	if err.Type != ErrNotFound {
		t.Errorf("Type = %v, want %v", err.Type, ErrNotFound)
	}
	if !strings.Contains(err.Message, "prod-us") {
		t.Errorf("Message should contain context name")
	}
}

func TestConsulAuthFailed(t *testing.T) {
	cause := errors.New("connection refused")
	err := ConsulAuthFailed("localhost:8500", cause)
	if err.Type != ErrAuth {
		t.Errorf("Type = %v, want %v", err.Type, ErrAuth)
	}
	if !strings.Contains(err.Message, "localhost:8500") {
		t.Errorf("Message should contain address")
	}
	if err.Cause != cause {
		t.Errorf("Cause = %v, want %v", err.Cause, cause)
	}
}

func TestAWSAuthFailed(t *testing.T) {
	cause := errors.New("invalid credentials")
	err := AWSAuthFailed("us-east-1", cause)
	if err.Type != ErrAuth {
		t.Errorf("Type = %v, want %v", err.Type, ErrAuth)
	}
	if !strings.Contains(err.Message, "us-east-1") {
		t.Errorf("Message should contain region")
	}
}

func TestProviderNotConfigured(t *testing.T) {
	err := ProviderNotConfigured("vault")
	if err.Type != ErrProvider {
		t.Errorf("Type = %v, want %v", err.Type, ErrProvider)
	}
	if !strings.Contains(err.Message, "vault") {
		t.Errorf("Message should contain provider name")
	}
}

func TestSecretNotFound(t *testing.T) {
	err := SecretNotFound("consul", "services/auth/password")
	if err.Type != ErrNotFound {
		t.Errorf("Type = %v, want %v", err.Type, ErrNotFound)
	}
	if !strings.Contains(err.Message, "services/auth/password") {
		t.Errorf("Message should contain key")
	}
	if !strings.Contains(err.Message, "consul") {
		t.Errorf("Message should contain provider")
	}
}

func TestCredentialResolutionFailed(t *testing.T) {
	cause := errors.New("provider error")
	err := CredentialResolutionFailed("auth-service", "prod", cause)
	if err.Type != ErrProvider {
		t.Errorf("Type = %v, want %v", err.Type, ErrProvider)
	}
	if !strings.Contains(err.Message, "auth-service") {
		t.Errorf("Message should contain service")
	}
	if !strings.Contains(err.Message, "prod") {
		t.Errorf("Message should contain env")
	}
}

func TestRequestFailed(t *testing.T) {
	cause := errors.New("timeout")
	err := RequestFailed("https://api.example.com/test", cause)
	if err.Type != ErrNetwork {
		t.Errorf("Type = %v, want %v", err.Type, ErrNetwork)
	}
	if !strings.Contains(err.Message, "https://api.example.com/test") {
		t.Errorf("Message should contain URL")
	}
}

func TestBaseURLMissing(t *testing.T) {
	err := BaseURLMissing("auth-service", "prod")
	if err.Type != ErrValidation {
		t.Errorf("Type = %v, want %v", err.Type, ErrValidation)
	}
}

func TestInvalidMethod(t *testing.T) {
	err := InvalidMethod("INVALID")
	if err.Type != ErrValidation {
		t.Errorf("Type = %v, want %v", err.Type, ErrValidation)
	}
	if !strings.Contains(err.Message, "INVALID") {
		t.Errorf("Message should contain method")
	}
}

func TestMissingRequiredFlag(t *testing.T) {
	err := MissingRequiredFlag("service")
	if err.Type != ErrValidation {
		t.Errorf("Type = %v, want %v", err.Type, ErrValidation)
	}
	if !strings.Contains(err.Message, "service") {
		t.Errorf("Message should contain flag name")
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("original")
	err := Wrap(cause, "wrapped message")
	if err.Message != "wrapped message" {
		t.Errorf("Message = %q, want %q", err.Message, "wrapped message")
	}
	if err.Cause != cause {
		t.Errorf("Cause = %v, want %v", err.Cause, cause)
	}
}

func TestProviderInitFailed(t *testing.T) {
	cause := errors.New("init error")
	err := ProviderInitFailed("Consul", cause)
	if err.Type != ErrProvider {
		t.Errorf("Type = %v, want %v", err.Type, ErrProvider)
	}
	if !strings.Contains(err.Message, "Consul") {
		t.Errorf("Message should contain provider name")
	}
}

func TestConsulAddressRequired(t *testing.T) {
	err := ConsulAddressRequired()
	if err.Type != ErrConfig {
		t.Errorf("Type = %v, want %v", err.Type, ErrConfig)
	}
}

func TestConsulKeyNotFound(t *testing.T) {
	err := ConsulKeyNotFound("services/auth/url")
	if err.Type != ErrNotFound {
		t.Errorf("Type = %v, want %v", err.Type, ErrNotFound)
	}
}

func TestConsulGetFailed(t *testing.T) {
	cause := errors.New("connection error")
	err := ConsulGetFailed("services/auth/url", cause)
	if err.Type != ErrProvider {
		t.Errorf("Type = %v, want %v", err.Type, ErrProvider)
	}
}

func TestServiceAlreadyExists(t *testing.T) {
	err := ServiceAlreadyExists("auth-service")
	if err.Type != ErrValidation {
		t.Errorf("Type = %v, want %v", err.Type, ErrValidation)
	}
}

func TestInvalidPathMapping(t *testing.T) {
	err := InvalidPathMapping("invalid-mapping")
	if err.Type != ErrValidation {
		t.Errorf("Type = %v, want %v", err.Type, ErrValidation)
	}
}

func TestServiceModeMixed(t *testing.T) {
	err := ServiceModeMixed()
	if err.Type != ErrValidation {
		t.Errorf("Type = %v, want %v", err.Type, ErrValidation)
	}
}

func TestServiceModeRequired(t *testing.T) {
	err := ServiceModeRequired()
	if err.Type != ErrValidation {
		t.Errorf("Type = %v, want %v", err.Type, ErrValidation)
	}
}

func TestPathResolutionFailed(t *testing.T) {
	cause := errors.New("path error")
	err := PathResolutionFailed("services/{service}/url", cause)
	if err.Type != ErrProvider {
		t.Errorf("Type = %v, want %v", err.Type, ErrProvider)
	}
}

func TestJSONKeyNotFound(t *testing.T) {
	err := JSONKeyNotFound("password", "secrets/db")
	if err.Type != ErrNotFound {
		t.Errorf("Type = %v, want %v", err.Type, ErrNotFound)
	}
}

func TestJSONParseFailed(t *testing.T) {
	cause := errors.New("json error")
	err := JSONParseFailed(cause)
	if err.Type != ErrValidation {
		t.Errorf("Type = %v, want %v", err.Type, ErrValidation)
	}
}

func TestErrorTypes(t *testing.T) {
	// Ensure all error types are distinct
	types := []ErrorType{ErrConfig, ErrAuth, ErrProvider, ErrNetwork, ErrNotFound, ErrValidation}
	seen := make(map[ErrorType]bool)
	for _, et := range types {
		if seen[et] {
			t.Errorf("Duplicate error type: %v", et)
		}
		seen[et] = true
	}
}
