package yaml

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DefaultFile      = "api.yaml"
	GoFireConfigFile = ".gofire.yaml"
)

type Endpoint struct {
	Method  string `yaml:"method"`
	Path    string `yaml:"path"`
	Handler string `yaml:"handler"`
	Auth    bool   `yaml:"auth"`
}

type OutputConfig struct {
	ServerDir   string `yaml:"serverDir,omitempty"`
	HandlersDir string `yaml:"handlersDir,omitempty"`
}

type APIConfig struct {
	Version   string        `yaml:"version"`
	BasePath  string        `yaml:"basePath"`
	Endpoints []Endpoint    `yaml:"endpoints"`
	Output    *OutputConfig `yaml:"output,omitempty"`
}

// GoFireConfig is the structure of .gofire.yaml (project-level overrides).
type GoFireConfig struct {
	Output *OutputConfig `yaml:"output,omitempty"`
}

// LoadGoFireConfig reads .gofire.yaml from the given path.
// Returns (nil, nil) if the file does not exist; returns error only when the file exists but fails to parse.
func LoadGoFireConfig(path string) (*GoFireConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	var cfg GoFireConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return &cfg, nil
}

// ResolveServerDir returns server directory. Priority: flag > .gofire.yaml > api.yaml output > default.
func ResolveServerDir(apiOutput, gofireOutput *OutputConfig, flagVal string) string {
	if flagVal != "" {
		return flagVal
	}
	if gofireOutput != nil && gofireOutput.ServerDir != "" {
		return gofireOutput.ServerDir
	}
	if apiOutput != nil && apiOutput.ServerDir != "" {
		return apiOutput.ServerDir
	}
	return "server"
}

// ResolveHandlersDir returns handlers directory. Priority: flag > .gofire.yaml > api.yaml output > default.
func ResolveHandlersDir(apiOutput, gofireOutput *OutputConfig, flagVal string) string {
	if flagVal != "" {
		return flagVal
	}
	if gofireOutput != nil && gofireOutput.HandlersDir != "" {
		return gofireOutput.HandlersDir
	}
	if apiOutput != nil && apiOutput.HandlersDir != "" {
		return apiOutput.HandlersDir
	}
	return "handlers"
}

// Load reads and parses the api.yaml file.
func Load(path string) (*APIConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var cfg APIConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	return &cfg, nil
}

// Save writes the APIConfig back to the yaml file.
func Save(path string, cfg *APIConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// AddEndpoint appends an endpoint to the config and saves.
func AddEndpoint(path string, method, route string, auth bool) error {
	cfg, err := Load(path)
	if err != nil {
		return err
	}

	handler := buildHandlerName(method, route)

	for _, ep := range cfg.Endpoints {
		if strings.EqualFold(ep.Method, method) && ep.Path == route {
			return fmt.Errorf("endpoint %s %s already exists", method, route)
		}
	}

	cfg.Endpoints = append(cfg.Endpoints, Endpoint{
		Method:  strings.ToUpper(method),
		Path:    route,
		Handler: handler,
		Auth:    auth,
	})

	return Save(path, cfg)
}

// buildHandlerName generates a handler function name from method + path.
// e.g. GET /users -> ListUsers, POST /users -> CreateUsers, GET /users/:id -> GetUsersById
func buildHandlerName(method, path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	hasParam := false
	var segments []string
	for _, p := range parts {
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, ":") {
			hasParam = true
			clean := strings.TrimPrefix(p, ":")
			segments = append(segments, "By"+capitalize(clean))
		} else {
			clean := strings.ReplaceAll(p, "-", "_")
			segments = append(segments, capitalize(clean))
		}
	}

	resource := strings.Join(segments, "")
	if resource == "" {
		resource = "Root"
	}

	var prefix string
	switch strings.ToUpper(method) {
	case "GET":
		if hasParam {
			prefix = "Get"
		} else {
			prefix = "List"
		}
	case "POST":
		prefix = "Create"
	case "PUT":
		prefix = "Update"
	case "PATCH":
		prefix = "Patch"
	case "DELETE":
		prefix = "Delete"
	default:
		prefix = "Handle"
	}

	return prefix + resource
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// DefaultConfig returns a fresh api.yaml with health endpoints.
func DefaultConfig() *APIConfig {
	return &APIConfig{
		Version:  "1",
		BasePath: "/api",
		Endpoints: []Endpoint{
			{Method: "GET", Path: "/api", Handler: "Health", Auth: false},
			{Method: "GET", Path: "/api/health", Handler: "Health", Auth: false},
		},
	}
}
