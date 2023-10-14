package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ServerCmd)
	rootCmd.AddCommand(VersionCmd)
	rootCmd.AddCommand(WatchCmd)
}

var rootCmd = &cobra.Command{
	Use:   "file-sync",
	Short: "File sync is a tool to sync files between two directories",
}

func Execute() error {
	return rootCmd.Execute()
}
