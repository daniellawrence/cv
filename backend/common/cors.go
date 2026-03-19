package common

import (
	"net/http"
	"os"
	"strings"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get CORS origins from environment variable(s)
		corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
		if corsOrigins == "" {
			// Fall back to single origin for backward compatibility
			corsOrigins = os.Getenv("CORS_ALLOWED_ORIGIN")
		}

		var allowedOrigins []string
		if corsOrigins != "" {
			allowedOrigins = strings.Split(corsOrigins, ",")
			for i := range allowedOrigins {
				allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
			}
		}

		// If no origins specified, use localhost default for development
		if len(allowedOrigins) == 0 {
			allowedOrigins = []string{"http://localhost"}
		}

		// Check if request origin is in allowed list
		requestOrigin := r.Header.Get("Origin")
		var allowed bool
		for _, origin := range allowedOrigins {
			if origin == requestOrigin || strings.HasSuffix(requestOrigin, "."+strings.TrimPrefix(origin, "http://")) || 
			   (strings.Contains(origin, "*") && matchWildcard(origin, requestOrigin)) {
				allowed = true
				break
			}
		}

		// Set the allowed origin header (echo back the actual request origin if allowed)
		if allowed && requestOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", requestOrigin)
		} else if len(allowedOrigins) > 0 {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigins[0])
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, traceparent, tracestate")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// matchWildcard checks if a wildcard pattern matches the request origin
func matchWildcard(pattern, origin string) bool {
	pattern = strings.TrimPrefix(pattern, "http://")
	pattern = strings.TrimPrefix(pattern, "https://")
	
	if strings.HasPrefix(pattern, "*.") {
		domain := strings.TrimPrefix(pattern, "*.")
		return strings.HasSuffix(origin, "."+domain) || origin == domain
	}
	
	return false
}
