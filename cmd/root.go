package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stackup-paymaster",
	Short: "ERC-4337 Paymaster",
	Long:  "A JSON-RPC server for ERC-4337 Verifying Paymasters.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
