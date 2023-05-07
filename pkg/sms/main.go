package sms

import (
	"cnc/pkg/utils"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/twilio/twilio-go/twiml"
)

type SMSDestinationType int64

const (
	Normal SMSDestinationType = iota
	Satellite
)

type SMSDestination struct {
	PhoneNumber string
	Type        SMSDestinationType
}

type SMSClient struct {
	Client      *twilio.RestClient
	PhoneNumber string
}

func InitClient(sid, key, sourceNumber string) *SMSClient {
	return &SMSClient{
		Client: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: sid,
			Password: key,
		}),
		PhoneNumber: sourceNumber,
	}
}

// Parse an outbound message to a format suited for text messaging. More than likely a list of
// messages, due to binning for the different destinations
func (client *SMSClient) serializeOutbound(msg string, dest *SMSDestination) []string {
	binSize := 500

	// If the destination is satellite, we need to limit the number of
	// chars to small batches so that there is a better chance the
	// messages complete
	if dest.Type == Satellite {
		binSize = utils.C.BIN_SIZE
	}

	outboundMessages := []string{}
	var end int
	for i := 0; i < len(msg); i += binSize {
		end = utils.Min(i+binSize, len(msg))
		outboundMessages = append(outboundMessages, msg[i:end])
	}

	return outboundMessages
}

// Send a message to a destination
func (client *SMSClient) SendMessage(fullMessage string, destination *SMSDestination) {
	var params *twilioApi.CreateMessageParams
	var msg string
	messages := client.serializeOutbound(fullMessage, destination)
	for i := 0; i < len(messages); i++ {
		msg = messages[i]

		params = &twilioApi.CreateMessageParams{}
		params.SetFrom(client.PhoneNumber)
		params.SetTo(destination.PhoneNumber)
		params.SetBody(msg)

		resp, err := client.Client.Api.CreateMessage(params)
		if err != nil {
			log.WithFields(log.Fields{
				"response": resp,
				"msg":      msg,
			}).Fatal(err.Error())
		} else {
			response, _ := json.Marshal(*resp)
			log.WithFields(log.Fields{
				"response": response,
				"msg":      msg,
			}).Debug("Successfully sent message")
			fmt.Printf("Sent message to %v\n", destination.PhoneNumber)
		}
	}
}

// Parse a message stream to TwiML
func (client *SMSClient) parseToTwiML(messages []string) string {
	elements := []twiml.Element{}
	var messagingMessage *twiml.MessagingMessage
	for i := 0; i < len(messages); i++ {
		messagingMessage = &twiml.MessagingMessage{Body: messages[i]}
		elements = append(elements, messagingMessage)
	}
	twimlMessages, err := twiml.Messages(elements)
	if err != nil {
		log.Fatal(err)
	}

	return twimlMessages
}
