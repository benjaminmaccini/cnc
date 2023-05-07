package cmd

import (
	"cnc/pkg/sms"
	"cnc/pkg/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var client *sms.SMSClient
var destination *sms.SMSDestination

// Args
var isSat bool
var destNumber string
var msg string
var sourceNumber string

// Root sms command
var smsCmd = &cobra.Command{
	Use:   "sms",
	Short: "Commands that handle the SMS entrypoint",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Prevent usage from an internal execution
		if internal {
			log.WithFields(log.Fields{
				"cmd":  cmd.CalledAs(),
				"args": args,
			}).Fatal("Tried to execute restricted command internally")
		}
		client = sms.InitClient(
			utils.C.TWILIO_ACCOUNT_SID,
			utils.C.TWILIO_KEY,
			sourceNumber,
		)

		if destNumber != "" {
			var destinationType sms.SMSDestinationType = sms.Normal
			if isSat {
				destinationType = sms.Satellite
			}

			destination = &sms.SMSDestination{
				PhoneNumber: destNumber,
				Type:        destinationType,
			}
		}
	},
}

var smsSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a message",
	Run: func(cmd *cobra.Command, args []string) {
		client.SendMessage(msg, destination)
	},
}

var smsReceiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "A receive a message",
	Long:  `This starts a web server with a webhook`,
	Run: func(cmd *cobra.Command, args []string) {
		sms.InitSMSServer(utils.C.SMS_SERVER_PORT, client, rootCmd)
	},
}

func init() {
	rootCmd.AddCommand(smsCmd)

	smsCmd.PersistentFlags().StringVarP(&destNumber, "destination", "d", "", "An external number to send to or receive from")
	smsCmd.PersistentFlags().StringVarP(&sourceNumber, "source", "s", "", "The client's assumed number")
	smsCmd.PersistentFlags().BoolVarP(&isSat, "satellite", "a", false, "If the destination is a satellite phone")

	// Send command
	smsCmd.AddCommand(smsSendCmd)
	smsSendCmd.Flags().StringVarP(&msg, "message", "m", "", "Message to be sent")
	smsSendCmd.MarkFlagRequired("message")
	smsSendCmd.MarkFlagRequired("destination")

	// Receive command
	smsCmd.AddCommand(smsReceiveCmd)
}
