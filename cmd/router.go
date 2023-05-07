package cmd

import (
	"cnc/pkg/router"
	"cnc/pkg/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var internalCommand string

// Root sms command
var routerCmd = &cobra.Command{
	Use:   "router",
	Short: "Convenience target for testing commands through the router",
}

var routerTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the router with a command",
	Run: func(cmd *cobra.Command, args []string) {
		response := router.RouteStringAsCommand(internalCommand, rootCmd)
		utils.L.L.WithFields(log.Fields{"response": response}).Info("Handled successfully via the router")
	},
}

func init() {
	rootCmd.AddCommand(routerCmd)
	routerCmd.PersistentFlags().StringVarP(&internalCommand, "cmd", "c", "", "A command")

	// Receive command
	routerCmd.AddCommand(routerTestCmd)
	smsSendCmd.MarkFlagRequired("cmd")
}
