package scaffold

const handlerTemplate = `package handlers

import (
	"encoding/json"
	"net/http"
{{if .HasParams}}
	"github.com/go-chi/chi/v5"
{{end}}
)

// {{.Name}} handles {{.Method}} {{.Path}}
func {{.Name}}(w http.ResponseWriter, r *http.Request) {
{{if .HasParams}}
{{.URLParamDoc}}
	{{.UseParams}}
{{end}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"handler": "{{.Name}}",
		"method":  "{{.Method}}",
		"path":    "{{.Path}}",
		"status":  "TODO: implement",
	})
}
`

const serverTemplate = `package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/mustafaaksoy/goFire/config"
	"github.com/mustafaaksoy/goFire/handlers"
	"github.com/mustafaaksoy/goFire/middleware"
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
{{range .PublicRoutes}}	r.{{.ChiMethod}}("{{.Path}}", handlers.{{.Handler}})
{{end}}
	// --- Auth-protected routes ---
{{if .HasAuthRoutes}}	if firebaseAuth != nil {
		r.Group(func(r chi.Router) {
			r.Use(firebaseAuth.Middleware)
{{range .AuthRoutes}}			r.{{.ChiMethod}}("{{.Path}}", handlers.{{.Handler}})
{{end}}		})
	} else {
{{range .AuthRoutes}}		r.{{.ChiMethod}}("{{.Path}}", handlers.{{.Handler}})
{{end}}	}
{{else}}	_ = firebaseAuth
{{end}}
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
`
