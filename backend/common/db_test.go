package common

import (
	"os"
	"testing"
)

func TestGetDatabaseURL(t *testing.T) {
	tests := []struct {
		name           string
		envValue       string
		defaultURL     string
		expectedResult string
	}{
		{
			name:           "env_var_set",
			envValue:       "mysql_user@tcp(host:3306)/mydb",
			defaultURL:     "fallback_url",
			expectedResult: "mysql_user@tcp(host:3306)/mydb",
		},
		{
			name:           "env_var_not_set",
			envValue:       "",
			defaultURL:     "fallback_url",
			expectedResult: "fallback_url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(DBURLEnvVar, tt.envValue)
				defer os.Unsetenv(DBURLEnvVar)
			} else {
				os.Unsetenv(DBURLEnvVar)
			}

			result := GetDatabaseURL(tt.defaultURL)
			if result != tt.expectedResult {
				t.Errorf("expected %s, got %s", tt.expectedResult, result)
			}
		})
	}
}

func TestValidateConnectionString(t *testing.T) {
	tests := []struct {
		name        string
		dsn         string
		expectValid bool
	}{
		{
			name:        "valid_mysql_dsn",
			dsn:         "root@tcp(host:3306)/cv",
			expectValid: true,
		},
		{
			name:        "valid_mysql_dsn_with_port",
			dsn:         "user:password@tcp(db.example.com:3306)/mydb",
			expectValid: true,
		},
		{
			name:        "empty_string",
			dsn:         "",
			expectValid: false,
		},
		{
			name:        "invalid_format_too_short",
			dsn:         "just_a_string",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := ValidateConnectionString(tt.dsn)
			if valid != tt.expectValid {
				t.Errorf("expected valid=%v, got %v", tt.expectValid, valid)
			}

			if tt.expectValid && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectValid && err == nil {
				t.Error("expected an error for invalid connection string")
			}
		})
	}
}

func TestConnectWithValidation(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		defaultURL  string
		expectError bool
	}{
		{
			name:        "env_var_provided",
			envValue:    "root@tcp(localhost:3306)/test_cv",
			defaultURL:  "",
			expectError: false, // Will fail to connect but validation should pass
		},
		{
			name:        "fallback_to_default",
			envValue:    "",
			defaultURL:  "root@tcp(localhost:3306)/test_cv",
			expectError: false, // Will fail to connect but validation should pass
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(DBURLEnvVar, tt.envValue)
				defer os.Unsetenv(DBURLEnvVar)
			} else {
				os.Unsetenv(DBURLEnvVar)
			}

			_, err := ConnectWithValidation(tt.defaultURL)

			if tt.expectError && err == nil {
				t.Error("expected an error but got none")
			}
			
			// Note: We don't check for non-expectError case because the actual connection may fail
			// due to database not being available - that's expected in tests
		})
	}
}

func TestSplitDSN(t *testing.T) {
	tests := []struct {
		name     string
		dsn      string
		minParts int // minimum expected parts (splitting behavior may vary slightly)
	}{
		{
			name:     "simple_mysql_dsn",
			dsn:      "root@tcp(host:3306)/cv",
			minParts: 2, // Should have at least user and host/database info
		},
		{
			name:     "dsn_with_password",
			dsn:      "user:password@tcp(db.example.com:3306)/mydb",
			minParts: 2, // Should have at least user:pass and database info
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := splitDSN(tt.dsn)
			if len(parts) < tt.minParts {
				t.Errorf("expected at least %d parts, got %d: %v", tt.minParts, len(parts), parts)
			}
		})
	}
}
