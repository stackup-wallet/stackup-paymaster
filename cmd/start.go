package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stackup-wallet/stackup-paymaster/internal/start"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a JSON-RPC server",
	Long:  "The start command will run a JSON-RPC server to enable processing of UserOperations for the connected Verifying Paymaster.",
	Run: func(cmd *cobra.Command, args []string) {
		start.Server()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
