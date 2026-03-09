package main

import (
	"fmt"
	"log"

	"github.com/messivite/goFire/config"
	"github.com/messivite/goFire/server"
)

func main() {
	const cyan = "\033[36m"
	const bold = "\033[1m"
	const reset = "\033[0m"

	fmt.Println()
	fmt.Println(cyan + "  ╔════════════════════════╗" + reset)
	fmt.Printf(cyan+"  ║  "+reset+bold+"GoFire"+reset+cyan+"  v"+config.Version+reset+cyan+"           ║"+reset+"\n")
	fmt.Println(cyan + "  ╚════════════════════════╝" + reset)
	fmt.Println()

	cfg := config.LoadFromEnv()

	if err := server.Run(cfg); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
