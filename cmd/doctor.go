package cmd

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/pterm/pterm"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/auth"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/config"
	GraphHelper "github.com/saricon83-sudo/Tyr365AdminCli/graphHelper"
	"github.com/spf13/cobra"
)

// doctorCmd represents the diagnostic check command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check the health and authentication status of all backing APIs and environment dependencies",
	Long: `Doctor diagnoses environment configuration issues, credentials, and API connection status.
It performs checks against:
- Local configuration keys
- Local Azure CLI status
- Microsoft Graph API (token generation)
- Teams Governance API (auth & connection)
- M365 Archiver API (auth & connection)
- Team Toolbox API (auth & connection)`,
	Run: func(cmd *cobra.Command, args []string) {
		pterm.DefaultHeader.Println("365Admin System Diagnostics (Doctor)")
		fmt.Println()

		// 1. Audit Local Config file keys
		pterm.DefaultSection.Println("1. Local Configuration File Audit")
		cfg := config.Get()
		requiredKeys := []string{
			"O365TenantName",
			"client_id",
			"client_secret",
			"resource",
			"archiverAddress",
			"archiverApp",
			"archiverSecret",
			"archiverResource",
			"teamToolBoxAdress",
			"365ManagementAppId",
			"365ManagementAppSecret",
			"teamToolboxAppId",
		}
		
		hasMissingKeys := false
		for _, key := range requiredKeys {
			if !cfg.IsSet(key) || cfg.GetString(key) == "" {
				pterm.Warning.Printf("Missing or empty configuration key: %s\n", key)
				hasMissingKeys = true
			}
		}
		if !hasMissingKeys {
			pterm.Success.Println("All required configuration keys are present in config.json")
		}
		fmt.Println()

		// 2. Audit Azure CLI
		pterm.DefaultSection.Println("2. Local Azure CLI Diagnostics")
		if _, err := exec.LookPath("az"); err != nil {
			pterm.Warning.Println("Azure CLI: 'az' executable not found in PATH")
		} else {
			pterm.Success.Println("Azure CLI: Command executable found in PATH")
			// Test if logged in
			authCmd := exec.Command("az", "account", "show", "--output", "none")
			if err := authCmd.Run(); err != nil {
				pterm.Warning.Println("Azure CLI: Installed but not logged in (run 'az login' first)")
			} else {
				pterm.Success.Println("Azure CLI: Logged in and active subscription detected")
			}
		}
		fmt.Println()

		// 3. Microsoft Graph API
		pterm.DefaultSection.Println("3. Microsoft Graph API Connectivity")
		gHelper := GraphHelper.NewGraphHelper()
		if err := gHelper.InitializeGraphForAppAuth(); err != nil {
			pterm.Error.Printf("Microsoft Graph API: Initialization failed: %v\n", err)
		} else {
			pterm.Success.Println("Microsoft Graph API: Client initialization succeeded")
			token, err := gHelper.GetAppToken()
			if err != nil {
				pterm.Error.Printf("Microsoft Graph API: Token acquisition failed: %v\n", err)
			} else if token == nil || *token == "" {
				pterm.Error.Println("Microsoft Graph API: Received an empty token")
			} else {
				pterm.Success.Println("Microsoft Graph API: Successfully generated OAuth2 client credentials token")
			}
		}
		fmt.Println()

		// 4. Teams Governance API
		pterm.DefaultSection.Println("4. Teams Governance API Diagnostics")
		govToken, err := auth.GetGovernanceToken()
		if err != nil {
			pterm.Error.Printf("Teams Governance API: Token acquisition failed: %v\n", err)
		} else if govToken == "" {
			pterm.Error.Println("Teams Governance API: Received an empty token")
		} else {
			pterm.Success.Println("Teams Governance API: Token acquisition successful")
			
			// Ping connection
			client := &http.Client{Timeout: 5 * time.Second}
			apiURL := cfg.GetString("resource") + "/api/teams/AverageProvisionTime"
			req, _ := http.NewRequest("GET", apiURL, nil)
			req.Header.Set("Authorization", "Bearer "+govToken)
			
			resp, err := client.Do(req)
			if err != nil {
				pterm.Error.Printf("Teams Governance API: Ping request failed: %v\n", err)
			} else {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					pterm.Success.Println("Teams Governance API: Network connection and endpoint ping successful (HTTP 200)")
				} else {
					pterm.Warning.Printf("Teams Governance API: Ping returned unexpected HTTP status: %d\n", resp.StatusCode)
				}
			}
		}
		fmt.Println()

		// 5. M365 Archiver API
		pterm.DefaultSection.Println("5. M365 Archiver API Diagnostics")
		archAddress := cfg.GetString("archiverAddress")
		if archAddress == "" {
			pterm.Warning.Println("M365 Archiver API: Skipping diagnostics (archiverAddress not set in configuration)")
		} else {
			archToken, err := auth.GetArchiverToken()
			if err != nil {
				pterm.Error.Printf("M365 Archiver API: Token acquisition failed: %v\n", err)
			} else if archToken == "" {
				pterm.Error.Println("M365 Archiver API: Received an empty token")
			} else {
				pterm.Success.Println("M365 Archiver API: Token acquisition successful")
				
				// Ping connection
				client := &http.Client{Timeout: 5 * time.Second}
				apiURL := archAddress + "/api/Archiver/GetJobs"
				req, _ := http.NewRequest("GET", apiURL, nil)
				req.Header.Set("Authorization", "Bearer "+archToken)
				
				resp, err := client.Do(req)
				if err != nil {
					pterm.Error.Printf("M365 Archiver API: Ping request failed: %v\n", err)
				} else {
					defer resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						pterm.Success.Println("M365 Archiver API: Network connection and endpoint ping successful (HTTP 200)")
					} else {
						pterm.Warning.Printf("M365 Archiver API: Ping returned unexpected HTTP status: %d\n", resp.StatusCode)
					}
				}
			}
		}
		fmt.Println()

		// 6. Team Toolbox API
		pterm.DefaultSection.Println("6. Team Toolbox API Diagnostics")
		toolboxAddress := cfg.GetString("teamToolBoxAdress")
		if toolboxAddress == "" {
			pterm.Warning.Println("Team Toolbox API: Skipping diagnostics (teamToolBoxAdress not set in configuration)")
		} else {
			toolboxToken, err := auth.GetTeamToolboxToken()
			if err != nil {
				pterm.Error.Printf("Team Toolbox API: Token acquisition failed: %v\n", err)
			} else if toolboxToken == "" {
				pterm.Error.Println("Team Toolbox API: Received an empty token")
			} else {
				pterm.Success.Println("Team Toolbox API: Token acquisition successful")
				
				// Ping connection
				client := &http.Client{Timeout: 5 * time.Second}
				apiURL := toolboxAddress + "/Tools/TestPolicy"
				req, _ := http.NewRequest("GET", apiURL, nil)
				req.Header.Set("Authorization", "Bearer "+toolboxToken)
				
				resp, err := client.Do(req)
				if err != nil {
					pterm.Error.Printf("Team Toolbox API: Ping request failed: %v\n", err)
				} else {
					defer resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						pterm.Success.Println("Team Toolbox API: Network connection and endpoint ping successful (HTTP 200)")
					} else {
						pterm.Warning.Printf("Team Toolbox API: Ping returned unexpected HTTP status: %d\n", resp.StatusCode)
					}
				}
			}
		}
		fmt.Println()
		
		pterm.DefaultSection.Println("Diagnostic audit finished.")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
