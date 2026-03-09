package scaffold

const registryTemplate = `package {{.PackageName}}

import "net/http"

var Registry = make(map[string]http.HandlerFunc)

func Register(name string, fn http.HandlerFunc) {
	Registry[name] = fn
}

func Get(name string) http.HandlerFunc {
	return Registry[name]
}
`

const healthHandlerTemplate = `package {{.PackageName}}

import (
	"encoding/json"
	"net/http"
)

func init() {
	Register("Health", Health)
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
`

const rootHandlerContent = "package {{.PackageName}}\n\nimport (\n\t\"net/http\"\n\n\t\"github.com/messivite/goFire/config\"\n)\n\nfunc init() {\n\tRegister(\"Root\", Root)\n}\n\nfunc Root(w http.ResponseWriter, r *http.Request) {\n\tw.Header().Set(\"Content-Type\", \"text/html; charset=utf-8\")\n\n\thtml := `<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n  <meta charset=\"UTF-8\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n  <title>GoFire</title>\n  <link href=\"https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@600;700&display=swap\" rel=\"stylesheet\">\n  <style>\n    * { margin: 0; padding: 0; box-sizing: border-box; }\n    body { min-height: 100vh; display: flex; align-items: center; justify-content: center; background: #0d1117; font-family: 'JetBrains Mono', monospace; color: #58a6ff; overflow: hidden; }\n    .container { text-align: center; animation: fadeIn 0.8s ease-out; }\n    @keyframes fadeIn { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }\n    .box { display: inline-block; padding: 2.5rem 3.5rem; border: 2px solid #58a6ff; border-radius: 6px; position: relative; }\n    .box::before { content: ''; position: absolute; inset: -1px; border-radius: 6px; background: linear-gradient(135deg, rgba(88,166,255,0.15), transparent); z-index: -1; }\n    .fire { font-size: 3rem; margin-bottom: 0.5rem; }\n    .brand { font-size: 2.8rem; font-weight: 700; color: #fff; letter-spacing: 0.15em; text-shadow: 0 0 30px rgba(88,166,255,0.4); }\n    .version { margin-top: 0.75rem; font-size: 0.85rem; color: #8b949e; }\n    .links { margin-top: 1.8rem; display: flex; gap: 1.5rem; justify-content: center; font-size: 0.8rem; }\n    .links a { color: #58a6ff; text-decoration: none; padding: 0.4rem 0.8rem; border: 1px solid #30363d; border-radius: 4px; transition: all 0.2s; }\n    .links a:hover { border-color: #58a6ff; background: rgba(88,166,255,0.1); }\n  </style>\n</head>\n<body>\n  <div class=\"container\">\n    <div class=\"box\">\n      <div class=\"brand\">GoFire</div>\n      <div class=\"version\">v` + config.Version + `</div>\n      <div class=\"links\">\n        <a href=\"/api/health\">/api/health</a>\n      </div>\n    </div>\n  </div>\n</body>\n</html>`\n\tw.Write([]byte(html))\n}\n"

const handlerTemplate = `package {{.PackageName}}

import (
	"encoding/json"
	"net/http"
{{if .HasParams}}
	"github.com/go-chi/chi/v5"
{{end}}
)

func init() {
	Register("{{.Name}}", {{.Name}})
}

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

	"github.com/messivite/goFire/apidef"
	"github.com/messivite/goFire/config"
	{{.HandlersPackage}} "{{.HandlersImportPath}}"
	"github.com/messivite/goFire/middleware"
)

func NewHandler(cfg *config.Config) (http.Handler, error) {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)

	apiCfg, err := apidef.Load(apidef.DefaultFile)
	if err != nil {
		return nil, fmt.Errorf("loading api.yaml: %w", err)
	}

	r.Get("/", {{.HandlersPackage}}.Root)

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

	for _, ep := range apiCfg.Endpoints {
		fn := {{.HandlersPackage}}.Get(ep.Handler)
		if fn == nil {
			log.Printf("WARNING: handler %q not registered, skipping %s %s", ep.Handler, ep.Method, ep.Path)
			continue
		}
		path := apidef.ToChiPath(ep.Path)
		if ep.Auth && firebaseAuth != nil {
			r.Group(func(r chi.Router) {
				r.Use(firebaseAuth.Middleware)
				r.Method(ep.Method, path, fn)
			})
		} else {
			r.Method(ep.Method, path, fn)
		}
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
`

const cmdMainTemplate = `package main

import (
	"log"

	"github.com/messivite/goFire/config"
	"{{.ModulePath}}/server"
)

func main() {
	cfg := config.LoadFromEnv()
	if err := server.Run(cfg); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
`
