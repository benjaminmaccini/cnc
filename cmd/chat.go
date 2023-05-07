package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"cnc/pkg/chat"
	"cnc/pkg/utils"
)

var query string

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Execute a query against ChatGPT",
	Run: func(cmd *cobra.Command, args []string) {
		client := chat.InitClient(utils.C.OPENAI_KEY)
		response := client.Query(query)

		log.WithFields(log.Fields{
			"query":    query,
			"response": response,
		}).Debug()

		fmt.Printf(response)
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().StringVarP(&query, "query", "q", "", "A query for ChatGPT")
}
