package apidef

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultFile = "api.yaml"

type Endpoint struct {
	Method  string `yaml:"method"`
	Path    string `yaml:"path"`
	Handler string `yaml:"handler"`
	Auth    bool   `yaml:"auth"`
}

type APIConfig struct {
	Version   string     `yaml:"version"`
	BasePath  string     `yaml:"basePath"`
	Endpoints []Endpoint `yaml:"endpoints"`
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

// ToChiPath converts /users/:id to /users/{id} for Chi routing.
func ToChiPath(path string) string {
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

// ChiMethod returns the Chi router method name (e.g. GET -> Get).
func ChiMethod(method string) string {
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
