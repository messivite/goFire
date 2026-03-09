package scaffold

import (
	"os"
	"path/filepath"
	"strings"
)

// ReadGoModModule returns the module path from go.mod, or empty string if not found.
func ReadGoModModule(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			mod := strings.TrimPrefix(line, "module ")
			return strings.TrimSpace(mod)
		}
	}
	return ""
}
