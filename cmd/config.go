package cmd

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/config"
	"github.com/spf13/cobra"
)

// configCmd represents the parent config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage and inspect CLI configurations",
	Long:  `View, validate, and check configuration values currently active in 365Admin.`,
}

// configShowCmd represents the config show subcommand
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a secure, structured tree visualization of active configurations",
	Long:  `Displays all active configuration parameters nested by resource type. Credentials and secrets are automatically masked for security.`,
	Run: func(cmd *cobra.Command, args []string) {
		pterm.DefaultHeader.Println("365Admin Configuration Inspector")
		fmt.Println()

		cfg := config.Get()

		// Construct config tree using pterm.TreeNode
		treeRoot := pterm.TreeNode{
			Text: pterm.LightMagenta("365Admin Configuration Registry"),
			Children: []pterm.TreeNode{
				{
					Text: pterm.LightCyan("Microsoft Graph Integration"),
					Children: []pterm.TreeNode{
						{Text: fmt.Sprintf("O365 Tenant Name: %s", formatVal(cfg.GetString("O365TenantName")))},
						{Text: fmt.Sprintf("Client ID:        %s", maskCredential(cfg.GetString("client_id")))},
						{Text: fmt.Sprintf("Client Secret:    %s", maskSecret(cfg.GetString("client_secret")))},
					},
				},
				{
					Text: pterm.LightCyan("Teams Governance API"),
					Children: []pterm.TreeNode{
						{Text: fmt.Sprintf("Resource Endpoint: %s", formatVal(cfg.GetString("resource")))},
						{Text: fmt.Sprintf("Auth Client ID:    %s", maskCredential(cfg.GetString("client_id")))},
						{Text: fmt.Sprintf("Auth Client Sec:   %s", maskSecret(cfg.GetString("client_secret")))},
						{Text: fmt.Sprintf("OAuth Grant Type:  %s", formatVal(cfg.GetString("grant_type")))},
					},
				},
				{
					Text: pterm.LightCyan("M365 Archiver Service"),
					Children: []pterm.TreeNode{
						{Text: fmt.Sprintf("Service Address:   %s", formatVal(cfg.GetString("archiverAddress")))},
						{Text: fmt.Sprintf("App Identifier:    %s", maskCredential(cfg.GetString("archiverApp")))},
						{Text: fmt.Sprintf("App Secret:        %s", maskSecret(cfg.GetString("archiverSecret")))},
						{Text: fmt.Sprintf("Resource Endpoint: %s", formatVal(cfg.GetString("archiverResource")))},
					},
				},
				{
					Text: pterm.LightCyan("Team Toolbox Toolset"),
					Children: []pterm.TreeNode{
						{Text: fmt.Sprintf("Service Address:   %s", formatVal(cfg.GetString("teamToolBoxAdress")))},
						{Text: fmt.Sprintf("Management App ID: %s", maskCredential(cfg.GetString("365ManagementAppId")))},
						{Text: fmt.Sprintf("Management Secret: %s", maskSecret(cfg.GetString("365ManagementAppSecret")))},
						{Text: fmt.Sprintf("Toolbox App ID:    %s", maskCredential(cfg.GetString("teamToolboxAppId")))},
					},
				},
			},
		}

		// Render the nested configuration tree
		err := pterm.DefaultTree.WithRoot(treeRoot).Render()
		if err != nil {
			pterm.Error.Printf("Failed to render configuration tree: %v\n", err)
		}
		
		fmt.Println()
		pterm.Info.Println("Tip: To check connection health for these endpoints, run '365Admin doctor'")
	},
}

func formatVal(val string) string {
	if val == "" {
		return pterm.Red("[NOT CONFIGURED]")
	}
	return pterm.Green(val)
}

func maskCredential(val string) string {
	if val == "" {
		return pterm.Red("[NOT CONFIGURED]")
	}
	if len(val) <= 8 {
		return pterm.Yellow("********")
	}
	return pterm.Green(val[:4] + "-****-****-****-" + val[len(val)-4:])
}

func maskSecret(val string) string {
	if val == "" {
		return pterm.Red("[NOT CONFIGURED]")
	}
	return pterm.Yellow("**************** (Masked)")
}

func init() {
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}
