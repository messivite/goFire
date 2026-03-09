package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type contextKey string

const UserContextKey contextKey = "firebaseUser"

type FirebaseAuth struct {
	client *auth.Client
}

// NewFirebaseAuth initializes Firebase Auth from a service account JSON file path.
func NewFirebaseAuth(credentialsPath string) (*FirebaseAuth, error) {
	ctx := context.Background()

	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("Firebase Auth middleware initialized")
	return &FirebaseAuth{client: client}, nil
}

// NewFirebaseAuthFromJSON initializes Firebase Auth from credentials JSON bytes.
// Use for Vercel/serverless where credentials live in an env var.
func NewFirebaseAuthFromJSON(jsonBytes []byte) (*FirebaseAuth, error) {
	ctx := context.Background()

	opt := option.WithCredentialsJSON(jsonBytes)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("Firebase Auth middleware initialized (from JSON)")
	return &FirebaseAuth{client: client}, nil
}

// Middleware returns an http.Handler middleware that verifies Firebase ID tokens.
func (fa *FirebaseAuth) Middleware(next http.Handler) http.Handler {
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

		token, err := fa.client.VerifyIDToken(r.Context(), parts[1])
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
