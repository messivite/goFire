package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/messivite/goFire/config"
	"github.com/messivite/goFire/handlers"
	"github.com/messivite/goFire/middleware"
)

func NewHandler(cfg *config.Config) (http.Handler, error) {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)

	r.Get("/", handlers.Root)

	var firebaseAuth *middleware.FirebaseAuth
	if cfg.FirebaseEnabled() {
		var err error
		if cfg.FirebaseCredentialsJSON != "" {
			firebaseAuth, err = middleware.NewFirebaseAuthFromJSON([]byte(cfg.FirebaseCredentialsJSON))
		} else {
			firebaseAuth, err = middleware.NewFirebaseAuth(cfg.FirebaseCredentialsPath)
		}
		if err != nil {
			return nil, fmt.Errorf("initializing Firebase auth: %w", err)
		}
	} else {
		log.Println("WARNING: Firebase auth is disabled. All routes are public.")
	}

	// --- Public routes ---
	r.Get("/api", handlers.Health)
	r.Get("/api/health", handlers.Health)

	// --- Auth-protected routes ---
	if firebaseAuth != nil {
		r.Group(func(r chi.Router) {
			r.Use(firebaseAuth.Middleware)
			r.Get("/users/{id}", handlers.GetUsersById)
		})
	} else {
		r.Get("/users/{id}", handlers.GetUsersById)
	}

	return r, nil
}

func Run(cfg *config.Config) error {
	h, err := NewHandler(cfg)
	if err != nil {
		return err
	}

	addr := ":" + cfg.Port
	log.Printf("GoFire server starting on %s", addr)
	return http.ListenAndServe(addr, h)
}
