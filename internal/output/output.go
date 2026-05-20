package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	currentFormat = "table"
	mu            sync.RWMutex
)

// SetFormat updates the global output format preference (table|json|yaml).
func SetFormat(format string) {
	mu.Lock()
	defer mu.Unlock()
	if format != "" {
		currentFormat = format
	}
}

// GetFormat retrieves the currently active output format.
func GetFormat() string {
	mu.RLock()
	defer mu.RUnlock()
	return currentFormat
}

// PrintResult formats and prints data to stdout based on the globally configured output format.
// Runs fallbackTablePrinter if output format is set to 'table'.
func PrintResult(data interface{}, fallbackTablePrinter func()) {
	format := GetFormat()

	switch format {
	case "json":
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting JSON output: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	case "yaml":
		yamlData, err := yaml.Marshal(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting YAML output: %v\n", err)
			return
		}
		fmt.Println(string(yamlData))
	default:
		if fallbackTablePrinter != nil {
			fallbackTablePrinter()
		} else {
			// JSON fallback
			jsonData, _ := json.MarshalIndent(data, "", "  ")
			fmt.Println(string(jsonData))
		}
	}
}
