package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yadunut/file-sync/internal/server"
	"github.com/yadunut/file-sync/internal/server/db"
	"github.com/yadunut/file-sync/internal/server/http"
	"github.com/yadunut/file-sync/internal/util"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the server.",
	Run: func(cmd *cobra.Command, args []string) {
		config := util.GetConfig()
		db := db.NewDB(config.DBPath)
		log := util.CreateLogger()
		server := http.NewHttpServer(server.CreateServer(db, log, config))
		server.Start()
	},
}
