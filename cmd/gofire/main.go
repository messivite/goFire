package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
  gofire init                                  Create default api.yaml and project files
  gofire setup                                 Interactive config (port, Firebase, Redis) and save to .env
  gofire add endpoint "METHOD /path" [--auth]  Add an endpoint
  gofire gen                                   Generate handlers + server from api.yaml
  gofire list                                  List all endpoints
  gofire deploy                                Interactive Vercel deploy`)
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
	fmt.Println("Run 'gofire gen' to generate handler and server files.")
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

	fmt.Println("Generating server...")
	if err := scaffold.GenerateServer(cfg, "server"); err != nil {
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
