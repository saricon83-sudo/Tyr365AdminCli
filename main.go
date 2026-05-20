package main

import (
	"fmt"
	"os"

	"github.com/saricon83-sudo/Tyr365AdminCli/cmd"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/config"
)

func main() {
	if err := config.Initialize("config.json"); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}
	fmt.Println()
	cmd.Execute()
}
