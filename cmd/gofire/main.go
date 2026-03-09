package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	apiyaml "github.com/messivite/goFire/internal/yaml"
	"github.com/messivite/goFire/internal/scaffold"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		cmdAdd()
	case "gen":
		cmdGen()
	case "list":
		cmdList()
	case "new":
		cmdNew()
	case "deploy":
		cmdDeploy()
	case "init":
		cmdInit()
	case "setup":
		cmdSetup()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`GoFire CLI

Usage:
  gofire new <name>                             Create new project from scratch
  gofire init                                   Create api.yaml and cmd/server in existing project
  gofire setup                                  Interactive config (port, Firebase, Redis) and save to .env
  gofire add endpoint "METHOD /path" [--auth]   Add an endpoint
  gofire gen                                    Generate handlers + server from api.yaml
  gofire list                                   List all endpoints
  gofire deploy                                 Interactive Vercel deploy`)
}

// --- setup ---

func cmdSetup() {
	reader := bufio.NewReader(os.Stdin)

	port := prompt(reader, "Server port", "8080")
	firebasePath := prompt(reader, "Firebase credentials JSON path (e.g. ./service-account.json - leave empty to skip auth)", "")

	useRedis := prompt(reader, "Enable Redis cache (Upstash)? (y/n)", "n")
	var redisURL, redisToken string
	if strings.ToLower(useRedis) == "y" {
		redisURL = prompt(reader, "Upstash Redis REST URL (e.g. https://your-db.upstash.io)", "")
		redisToken = prompt(reader, "Upstash Redis REST Token", "")
	}

	save := prompt(reader, "Save configuration to .env file? (y/n)", "n")
	if strings.ToLower(save) != "y" && strings.ToLower(save) != "yes" {
		fmt.Println("Configuration not saved.")
		return
	}

	content := fmt.Sprintf("PORT=%s\nFIREBASE_CREDENTIALS_PATH=%s\n",
		port, firebasePath)
	if redisURL != "" || redisToken != "" {
		content += fmt.Sprintf("UPSTASH_REDIS_REST_URL=%s\nUPSTASH_REDIS_REST_TOKEN=%s\n",
			redisURL, redisToken)
	}

	if err := os.WriteFile(".env", []byte(content), 0600); err != nil {
		fmt.Printf("Error saving .env: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Configuration saved to .env")
}

func prompt(reader *bufio.Reader, label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("%s: ", label)
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}

// --- new ---

func cmdNew() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: gofire new <project-name>")
		fmt.Println("Example: gofire new my-api")
		os.Exit(1)
	}
	name := strings.TrimSpace(os.Args[2])
	if name == "" {
		fmt.Println("Project name cannot be empty.")
		os.Exit(1)
	}
	if strings.ContainsAny(name, `/\:*?"<>|`) {
		fmt.Println("Project name contains invalid characters.")
		os.Exit(1)
	}

	if _, err := os.Stat(name); err == nil {
		fmt.Printf("Directory %q already exists. Choose a different name or remove it.\n", name)
		os.Exit(1)
	}

	if err := os.MkdirAll(name, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Creating project %q...\n", name)

	// go mod init
	cmd := exec.Command("go", "mod", "init", name)
	cmd.Dir = name
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running go mod init: %v\n", err)
		os.Exit(1)
	}

	// go get goFire
	cmd = exec.Command("go", "get", "github.com/messivite/goFire")
	cmd.Dir = name
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running go get: %v\n", err)
		os.Exit(1)
	}

	// api.yaml
	cfg := apiyaml.DefaultConfig()
	yamlPath := filepath.Join(name, apiyaml.DefaultFile)
	if err := apiyaml.Save(yamlPath, cfg); err != nil {
		fmt.Printf("Error creating api.yaml: %v\n", err)
		os.Exit(1)
	}

	// cmd/server/main.go
	if err := scaffold.GenerateCmdServer(filepath.Join(name, "cmd", "server", "main.go"), name); err != nil {
		fmt.Printf("Error creating main.go: %v\n", err)
		os.Exit(1)
	}

	// handlers + server
	if err := scaffold.GenerateHandlers(cfg, filepath.Join(name, "handlers")); err != nil {
		fmt.Printf("Error generating handlers: %v\n", err)
		os.Exit(1)
	}
	if err := scaffold.GenerateServer(cfg, filepath.Join(name, "server"), name); err != nil {
		fmt.Printf("Error generating server: %v\n", err)
		os.Exit(1)
	}

	// .gitignore
	gitignore := ".env\nservice-account*.json\n"
	if err := os.WriteFile(filepath.Join(name, ".gitignore"), []byte(gitignore), 0644); err != nil {
		// non-fatal
	}

	// Makefile
	makefile := "run:\n\tgo run ./cmd/server\n\nbuild:\n\tgo build -o bin/server ./cmd/server\n"
	if err := os.WriteFile(filepath.Join(name, "Makefile"), []byte(makefile), 0644); err != nil {
		// non-fatal
	}

	fmt.Printf("\nCreated project %q.\n", name)
	fmt.Printf("  cd %s\n", name)
	fmt.Println("  go mod tidy")
	fmt.Println("  make run")
}

// --- init ---

func cmdInit() {
	yamlPath := apiyaml.DefaultFile

	if _, err := os.Stat(yamlPath); err == nil {
		fmt.Printf("%s already exists. Skipping.\n", yamlPath)
		return
	}

	cfg := apiyaml.DefaultConfig()
	if err := apiyaml.Save(yamlPath, cfg); err != nil {
		fmt.Printf("Error creating %s: %v\n", yamlPath, err)
		os.Exit(1)
	}

	fmt.Printf("Created %s with default health endpoints.\n", yamlPath)

	mainPath := filepath.Join("cmd", "server", "main.go")
	if _, err := os.Stat(mainPath); err == nil {
		fmt.Printf("%s already exists. Skipping.\n", mainPath)
	} else {
		wd, _ := os.Getwd()
		modulePath := scaffold.ReadGoModModule(wd)
		if modulePath == "" {
			modulePath = "example"
			fmt.Println("WARNING: go.mod not found. Using module path 'example' — run 'go mod init <module>' and fix cmd/server/main.go import if needed.")
		}
		if err := scaffold.GenerateCmdServer(mainPath, modulePath); err != nil {
			fmt.Printf("Error creating %s: %v\n", mainPath, err)
			os.Exit(1)
		}
		fmt.Printf("Created %s\n", mainPath)
	}

	fmt.Println("Run 'gofire gen' to generate handlers and server, then 'go run ./cmd/server'.")
}

// --- add ---

func cmdAdd() {
	if len(os.Args) < 4 || os.Args[2] != "endpoint" {
		fmt.Println("Usage: gofire add endpoint \"METHOD /path\" [--auth]")
		os.Exit(1)
	}

	spec := os.Args[3]
	auth := false
	for _, arg := range os.Args[4:] {
		if arg == "--auth" {
			auth = true
		}
	}

	parts := strings.Fields(spec)
	if len(parts) != 2 {
		fmt.Println("Invalid endpoint spec. Use: \"METHOD /path\" (e.g. \"GET /users\")")
		os.Exit(1)
	}

	method := strings.ToUpper(parts[0])
	path := parts[1]

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	yamlPath := apiyaml.DefaultFile
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		fmt.Printf("%s not found. Run 'gofire init' first.\n", yamlPath)
		os.Exit(1)
	}

	if err := apiyaml.AddEndpoint(yamlPath, method, path, auth); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	authLabel := ""
	if auth {
		authLabel = " [auth]"
	}
	fmt.Printf("Added: %s %s%s\n", method, path, authLabel)
	fmt.Println("Run 'gofire gen' to regenerate handlers and server.")
}

// --- gen ---

func cmdGen() {
	yamlPath := apiyaml.DefaultFile

	cfg, err := apiyaml.Load(yamlPath)
	if err != nil {
		fmt.Printf("Error loading %s: %v\n", yamlPath, err)
		os.Exit(1)
	}

	fmt.Println("Generating handlers...")
	if err := scaffold.GenerateHandlers(cfg, "handlers"); err != nil {
		fmt.Printf("Error generating handlers: %v\n", err)
		os.Exit(1)
	}

	wd, _ := os.Getwd()
	modulePath := scaffold.ReadGoModModule(wd)
	if modulePath == "" {
		modulePath = "example"
	}

	fmt.Println("Generating server...")
	if err := scaffold.GenerateServer(cfg, "server", modulePath); err != nil {
		fmt.Printf("Error generating server: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done! Run 'go build ./...' to verify.")
}

// --- list ---

func cmdList() {
	yamlPath := apiyaml.DefaultFile

	cfg, err := apiyaml.Load(yamlPath)
	if err != nil {
		fmt.Printf("Error loading %s: %v\n", yamlPath, err)
		os.Exit(1)
	}

	fmt.Printf("\nEndpoints (basePath: %s)\n", cfg.BasePath)
	fmt.Println(strings.Repeat("-", 50))

	for _, ep := range cfg.Endpoints {
		authLabel := "      "
		if ep.Auth {
			authLabel = "[auth]"
		}
		fmt.Printf("  %-7s %-25s %s  → %s\n", ep.Method, ep.Path, authLabel, ep.Handler)
	}
	fmt.Println()
}

// --- deploy ---

func cmdDeploy() {
	if _, err := exec.LookPath("vercel"); err != nil {
		fmt.Println("Error: 'vercel' CLI not found. Install it: npm i -g vercel")
		os.Exit(1)
	}

	if _, err := os.Stat(".vercel/project.json"); os.IsNotExist(err) {
		fmt.Println("No Vercel project linked. Running 'vercel link'...")
		runInteractive("vercel", "link")
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Deploy to production? (y/n) [n]: ")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "y" || answer == "yes" {
		fmt.Println("Deploying to production...")
		runInteractive("vercel", "--prod")
	} else {
		fmt.Println("Deploying preview...")
		runInteractive("vercel")
	}
}

func runInteractive(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Command failed: %v\n", err)
		os.Exit(1)
	}
}
