package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yadunut/file-sync/internal/client"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "gets the version",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewClient()
		c.Version()
	},
}
