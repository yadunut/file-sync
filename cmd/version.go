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
		log := util.CreateLogger()
		c := client.NewClient(log, util.GetConfig())
		v, err := c.Version()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(v.Version)
	},
}
