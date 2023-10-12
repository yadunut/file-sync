package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yadunut/file-sync/internal/client"
	"github.com/yadunut/file-sync/internal/util"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "gets the version",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewClient(util.GetConfig())
		v := c.Version()
		fmt.Println(v.Version)
	},
}
