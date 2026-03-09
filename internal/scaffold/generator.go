package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	apiyaml "github.com/messivite/goFire/internal/yaml"
)

type handlerData struct {
	PackageName string
	Name        string
	Method      string
	Path        string
	HasParams   bool
	URLParamDoc string
	UseParams   string // e.g. "_ = id" to avoid unused var
}

type routeData struct {
	ChiMethod string
	Path      string // Chi format: :id -> {id}
	Handler   string
}

// toChiPath converts /users/:id to /users/{id} for Chi routing.
func toChiPath(path string) string {
	var b strings.Builder
	i := 0
	for i < len(path) {
		if path[i] == ':' && i+1 < len(path) {
			b.WriteByte('{')
			i++
			for i < len(path) && isPathParamChar(path[i]) {
				b.WriteByte(path[i])
				i++
			}
			b.WriteByte('}')
		} else {
			b.WriteByte(path[i])
			i++
		}
	}
	return b.String()
}

func isPathParamChar(c byte) bool {
	return c == '_' || c == '-' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

// paramInfo holds original (for URLParam) and var (for Go) names
type paramInfo struct {
	Key string // chi.URLParam(r, key)
	Var string // Go variable name
}

func extractParamNames(path string) []paramInfo {
	var params []paramInfo
	parts := strings.Split(path, "/")
	for _, p := range parts {
		if strings.HasPrefix(p, ":") && len(p) > 1 {
			key := p[1:]
			params = append(params, paramInfo{Key: key, Var: strings.ReplaceAll(key, "-", "_")})
		}
	}
	return params
}

type serverData struct {
	ModulePath         string
	HandlersPackage    string
	HandlersImportPath string
	PublicRoutes       []routeData
	AuthRoutes         []routeData
	HasAuthRoutes      bool
}

// GenerateHandlers creates handler stub files for endpoints that don't already have one.
// Built-in handlers (Health, Root) get specific implementations; others get generic stubs.
// The Go package name is derived from filepath.Base(handlersDir).
func GenerateHandlers(cfg *apiyaml.APIConfig, handlersDir string) error {
	if err := os.MkdirAll(handlersDir, 0755); err != nil {
		return err
	}

	pkgName := filepath.Base(handlersDir)

	registryPath := filepath.Join(handlersDir, "registry.go")
	if _, err := os.Stat(registryPath); err != nil {
		tmpl, err := template.New("registry").Parse(registryTemplate)
		if err != nil {
			return err
		}
		f, err := os.Create(registryPath)
		if err != nil {
			return err
		}
		err = tmpl.Execute(f, struct{ PackageName string }{PackageName: pkgName})
		f.Close()
		if err != nil {
			return err
		}
		fmt.Printf("  created %s\n", registryPath)
	}

	builtinTemplates := map[string]string{
		"Health": healthHandlerTemplate,
		"Root":   rootHandlerContent,
	}

	for _, name := range []string{"Health", "Root"} {
		filename := strings.ToLower(name) + ".go"
		fullPath := filepath.Join(handlersDir, filename)
		if _, err := os.Stat(fullPath); err == nil {
			continue
		}
		tmplStr, ok := builtinTemplates[name]
		if !ok {
			continue
		}
		tmpl, err := template.New(name).Parse(tmplStr)
		if err != nil {
			return err
		}
		f, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		err = tmpl.Execute(f, struct{ PackageName string }{PackageName: pkgName})
		f.Close()
		if err != nil {
			return err
		}
		fmt.Printf("  created %s\n", fullPath)
	}

	builtins := map[string]bool{"Health": true, "Root": true}
	seen := map[string]bool{}

	for _, ep := range cfg.Endpoints {
		if builtins[ep.Handler] || seen[ep.Handler] {
			continue
		}
		seen[ep.Handler] = true

		filename := strings.ToLower(ep.Handler) + ".go"
		fullPath := filepath.Join(handlersDir, filename)

		if _, err := os.Stat(fullPath); err == nil {
			fmt.Printf("  skip %s (already exists)\n", fullPath)
			continue
		}

		tmpl, err := template.New("handler").Parse(handlerTemplate)
		if err != nil {
			return err
		}

		f, err := os.Create(fullPath)
		if err != nil {
			return err
		}

		params := extractParamNames(ep.Path)
		var paramLines []string
		var useLines []string
		for _, p := range params {
			paramLines = append(paramLines, fmt.Sprintf("\t%s := chi.URLParam(r, %q)", p.Var, p.Key))
			useLines = append(useLines, fmt.Sprintf("_ = %s", p.Var))
		}
		paramDoc := strings.Join(paramLines, "\n")
		useParams := strings.Join(useLines, "\n\t")

		err = tmpl.Execute(f, handlerData{
			PackageName: pkgName,
			Name:        ep.Handler,
			Method:      ep.Method,
			Path:        ep.Path,
			HasParams:   len(params) > 0,
			URLParamDoc: paramDoc,
			UseParams:   useParams,
		})
		f.Close()
		if err != nil {
			return err
		}

		fmt.Printf("  created %s\n", fullPath)
	}

	return nil
}

// GenerateServer writes server/server.go from the api.yaml config.
// modulePath is read from go.mod; if empty, falls back to "example".
// handlersDir is the relative path to the handlers directory (e.g. "handlers", "pkg/handler").
func GenerateServer(cfg *apiyaml.APIConfig, serverDir string, modulePath string, handlersDir string) error {
	if err := os.MkdirAll(serverDir, 0755); err != nil {
		return err
	}
	if modulePath == "" {
		modulePath = "example"
	}
	if handlersDir == "" {
		handlersDir = "handlers"
	}

	handlersPackage := filepath.Base(handlersDir)
	handlersImportPath := modulePath + "/" + filepath.ToSlash(handlersDir)

	var sd serverData
	sd.ModulePath = modulePath
	sd.HandlersPackage = handlersPackage
	sd.HandlersImportPath = handlersImportPath
	for _, ep := range cfg.Endpoints {
		rd := routeData{
			ChiMethod: chiMethod(ep.Method),
			Path:      toChiPath(ep.Path),
			Handler:   ep.Handler,
		}
		if ep.Auth {
			sd.AuthRoutes = append(sd.AuthRoutes, rd)
		} else {
			sd.PublicRoutes = append(sd.PublicRoutes, rd)
		}
	}
	sd.HasAuthRoutes = len(sd.AuthRoutes) > 0

	tmpl, err := template.New("server").Parse(serverTemplate)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(serverDir, "server.go")
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := tmpl.Execute(f, sd); err != nil {
		return err
	}

	fmt.Printf("  generated %s\n", fullPath)
	return nil
}

// GenerateCmdServer creates cmd/server/main.go. modulePath from go.mod; if empty, uses "example".
func GenerateCmdServer(mainPath string, modulePath string) error {
	dir := filepath.Dir(mainPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if modulePath == "" {
		modulePath = "example"
	}
	tmpl, err := template.New("cmdmain").Parse(cmdMainTemplate)
	if err != nil {
		return err
	}
	f, err := os.Create(mainPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, struct{ ModulePath string }{ModulePath: modulePath})
}

func chiMethod(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return "Get"
	case "POST":
		return "Post"
	case "PUT":
		return "Put"
	case "PATCH":
		return "Patch"
	case "DELETE":
		return "Delete"
	default:
		return "Method"
	}
}
