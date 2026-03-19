package common

import (
	"net/http"
	"strconv"
	"strings"
)

func GetPathInt(r *http.Request, name string, defaultValue int) int {
	v := r.PathValue(name)
	if v == "" {
		return defaultValue
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return n
}

// ValidateID validates that an ID parameter is in a safe format (numeric only)
func ValidateID(id string) bool {
	if id == "" {
		return false
	}
	// Ensure ID contains only digits and no special characters
	for _, c := range id {
		if c < '0' || c > '9' {
			return false
		}
	}
	// Reject IDs that are too long (potential overflow attack)
	if len(id) > 15 {
		return false
	}
	// Reject leading zeros to prevent ambiguity
	if len(id) > 1 && id[0] == '0' {
		return false
	}
	return true
}

// ValidateIDOrEmpty allows empty IDs (for cases where empty is valid)
func ValidateIDOrEmpty(id string) bool {
	if id == "" {
		return true
	}
	return ValidateID(id)
}

// ValidateStringID validates that an ID parameter contains only safe characters
// Allows alphanumeric characters, hyphens, and underscores
func ValidateStringID(id string) bool {
	if id == "" {
		return false
	}
	// Ensure ID contains only alphanumeric characters, hyphens, and underscores
	for _, c := range id {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || 
		     (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}
	// Reject IDs that are too long (potential overflow attack)
	if len(id) > 100 {
		return false
	}
	return true
}

// ValidateStringIDOrEmpty allows empty IDs (for cases where empty is valid)
func ValidateStringIDOrEmpty(id string) bool {
	if id == "" {
		return true
	}
	return ValidateStringID(id)
}

// SanitizeString removes potentially dangerous characters from user input
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}
