package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yadunut/file-sync/internal/client"
	"github.com/yadunut/file-sync/internal/contracts"
	"github.com/yadunut/file-sync/internal/util"
)

func init() {
	WatchCmd.AddCommand(WatchUpCmd)
	WatchCmd.AddCommand(WatchDownCmd)
	WatchCmd.AddCommand(WatchListCmd)
}

var WatchCmd = &cobra.Command{
	Use: "watch",
}

var log = util.CreateLogger()
var c = client.NewClient(log, util.GetConfig())

var WatchUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Adds a directory to be watched",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path, err := filepath.Abs(args[0])
		if err != nil {
			log.Fatal(err)
		}
		res, err := c.WatchUp(contracts.WatchUpReq{Path: path})
		if res.Success {
			log.Info("Success")
		} else {
			log.Info("Failure")
		}
	},
}
var WatchDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Removes a directory from being watched",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path, err := filepath.Abs(args[0])
		if err != nil {
			log.Fatal(err)
		}
		res, err := c.WatchDown(contracts.WatchDownReq{Path: path})
		if err != nil {
			log.Fatal(err)
		}
		if res.Success {
			log.Info("Success")
		} else {
			log.Info("Failure")
		}
	},
}
var WatchListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all directories being watched",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := c.WatchList()
		if err != nil {
			log.Fatal(err)
		}
		if !res.Success {
			log.Error("Failure")
			return
		}
		log.Info("Success")
		if len(res.Directories) == 0 {
			log.Info("No files currently being watched")
			return
		}
		for _, dir := range res.Directories {
			log.Info(dir.Path)
		}
	},
}
