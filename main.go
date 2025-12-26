package main

import (
	"fmt"

	"github.com/saricon83-sudo/Tyr365AdminCli/cmd"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/config"
)

func main() {

	if err := config.Initialize("config.json"); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	fmt.Println()
	cmd.Execute()
}
