package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const demoToken = "demo-token"

// DemoAuth returns middleware for demo mode: only "Bearer demo-token" is accepted.
// Set GOFIRE_DEMO_AUTH=1 to enable. No Firebase required.
func DemoAuth() func(http.Handler) http.Handler {
	log.Println("Demo auth mode enabled. Use 'Authorization: Bearer demo-token' for protected routes.")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing Authorization header"})
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid Authorization format, expected: Bearer <token>"})
				return
			}

			if parts[1] != demoToken {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
