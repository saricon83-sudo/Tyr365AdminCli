package teamGov

// import (
//     "fmt"

//     teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
//     logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
//     log "github.com/sirupsen/logrus"
//     "github.com/spf13/cobra"
// )



// // lockTeamCmd represents the lockTeam command
// var lockTeamCmd = &cobra.Command{
//     Use:   "lockTeam",
//     Short: "Locks a team to prevent modifications",
//     Long: `Locks a team in the Teams Governance API using the provided groupId.
// This prevents modifications to the team until it is unlocked.

// Example:
//   365Admin teamGov lockTeam --groupId "abc123-def456-ghi789"`,
//     Run: func(cmd *cobra.Command, args []string) {
//         logger := logging.GetLogger()

//         if cmd.Flag("help").Changed {
//             cmd.Help()
//             return
//         }

//         body, err := teamGovHttp.GetFlexible("/api/Group/LockTeam", map[string]string{"groupId": groupId})
//         if err != nil {
//             logger.WithFields(log.Fields{
//                 "url":    "/api/Group/LockTeam",
//                 "method": "GET",
//                 "status": "Error",
//             }).Error(err)
//             return
//         }

//         fmt.Println(string(body))
//     },
// }

// func init() {
//     lockTeamCmd.Flags().StringVarP(&groupId, "groupId", "g", "", "The groupId of the team to lock")
//     if err := lockTeamCmd.MarkFlagRequired("groupId"); err != nil {
//         fmt.Println(err)
//     }

//     TeamGovCmd.AddCommand(lockTeamCmd)
// }