package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/yadunut/file-sync/internal/server"
	"github.com/yadunut/file-sync/internal/server/db"
	"github.com/yadunut/file-sync/internal/util"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the server.",
	Run: func(cmd *cobra.Command, args []string) {
		db := db.NewDB("./test.db")
		server := server.CreateServer(db, log.Default(), util.GetConfig())
		server.Start()
	},
}
