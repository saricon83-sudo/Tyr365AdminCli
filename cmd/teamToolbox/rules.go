package teamToolboxCmd

import (
	"fmt"
	"strconv"

	"github.com/pterm/pterm"
	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/output"
	"github.com/spf13/cobra"
)

var (
	rlId     int
	rlToolId int
)

// rulesCmd represents the rules command group
var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "Manage and inspect tool rules and logic",
}

// requirementsCmd represents the requirements command group
var requirementsCmd = &cobra.Command{
	Use:   "requirements",
	Short: "Inspect extended tool requirements",
}

// listRulesCmd represents the list command under rules
var listRulesCmd = &cobra.Command{
	Use:   "list",
	Short: "Get all registered rules",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching rules...")
		res, err := adminAPI.GetRules()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		output.PrintResult(res, func() { printRules(res) })
	},
}

// getRuleCmd represents the get command under rules
var getRuleCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a specific rule by ID",
	Run: func(cmd *cobra.Command, args []string) {
		if rlId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Rule ID")
			if err != nil {
				pterm.Error.Println("Rule ID is required")
				return
			}
			rlId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching rule ID %d...", rlId))
		rule, err := adminAPI.GetRule(rlId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		output.PrintResult([]teamToolboxHelper.TblToolRule{*rule}, func() {
			printRules([]teamToolboxHelper.TblToolRule{*rule})
		})
	},
}

// unusedRulesCmd represents the unused command under rules
var unusedRulesCmd = &cobra.Command{
	Use:   "unused",
	Short: "Get all unused rules",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching unused rules...")
		res, err := adminAPI.GetUnusedRules()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		output.PrintResult(res, func() { printRules(res) })
	},
}

// getRuleLogicsCmd represents the logics command under rules
var getRuleLogicsCmd = &cobra.Command{
	Use:   "logics",
	Short: "Get all rule logics",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching rule logics...")
		res, err := adminAPI.GetRuleLogics()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		output.PrintResult(res, func() { printLogics(res) })
	},
}

// getRuleLogicsByToolCmd represents the logics-by-tool command under rules
var getRuleLogicsByToolCmd = &cobra.Command{
	Use:   "logics-by-tool",
	Short: "Get rule logics for a specific tool",
	Run: func(cmd *cobra.Command, args []string) {
		if rlToolId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Tool ID")
			if err != nil {
				pterm.Error.Println("Tool ID is required")
				return
			}
			rlToolId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching rule logics for tool ID %d...", rlToolId))
		res, err := adminAPI.GetRuleLogicsByTool(rlToolId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		output.PrintResult(res, func() { printLogics(res) })
	},
}

// getExtendedRequirementsByToolCmd represents the by-tool command under requirements
var getExtendedRequirementsByToolCmd = &cobra.Command{
	Use:   "by-tool",
	Short: "Get extended requirements for a specific tool",
	Run: func(cmd *cobra.Command, args []string) {
		if rlToolId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Tool ID")
			if err != nil {
				pterm.Error.Println("Tool ID is required")
				return
			}
			rlToolId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching extended requirements for tool ID %d...", rlToolId))
		res, err := adminAPI.GetExtendedRequirementsByTool(rlToolId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		output.PrintResult(res, func() {
			tableData := pterm.TableData{
				{"ID", "Tool ID", "Requirement Name", "Value"},
			}
			for _, r := range res {
				tableData = append(tableData, []string{
					strconv.Itoa(r.Id),
					strconv.Itoa(r.ToolId),
					r.RequirementName,
					r.RequirementValue,
				})
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

// Helper printers
func printRules(rules []teamToolboxHelper.TblToolRule) {
	if len(rules) == 0 {
		pterm.Info.Println("No rules found.")
		return
	}
	tableData := pterm.TableData{
		{"ID", "Rule Name"},
	}
	for _, r := range rules {
		tableData = append(tableData, []string{
			strconv.Itoa(r.Id),
			r.RuleName,
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func printLogics(logics []teamToolboxHelper.TblToolRuleLogic) {
	if len(logics) == 0 {
		pterm.Info.Println("No rule logics found.")
		return
	}
	tableData := pterm.TableData{
		{"ID", "Rule Name", "Tool ID", "Rule ID", "Value", "Logic"},
	}
	for _, l := range logics {
		tableData = append(tableData, []string{
			strconv.Itoa(l.Id),
			l.RuleName,
			strconv.Itoa(l.ToolId),
			strconv.Itoa(l.RuleId),
			l.Value,
			l.Logic,
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func init() {
	getRuleCmd.Flags().IntVar(&rlId, "id", 0, "Rule ID")
	getRuleLogicsByToolCmd.Flags().IntVar(&rlToolId, "id", 0, "Tool ID")
	getExtendedRequirementsByToolCmd.Flags().IntVar(&rlToolId, "id", 0, "Tool ID")

	rulesCmd.AddCommand(listRulesCmd)
	rulesCmd.AddCommand(getRuleCmd)
	rulesCmd.AddCommand(unusedRulesCmd)
	rulesCmd.AddCommand(getRuleLogicsCmd)
	rulesCmd.AddCommand(getRuleLogicsByToolCmd)

	requirementsCmd.AddCommand(getExtendedRequirementsByToolCmd)

	TeamToolboxCmd.AddCommand(rulesCmd)
	TeamToolboxCmd.AddCommand(requirementsCmd)
}
