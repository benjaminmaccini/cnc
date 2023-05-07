package sms

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"cnc/pkg/router"
	"cnc/pkg/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type TwilioWebhookRequestBody struct {
	ToCountry     string
	ToState       string
	SmsMessageSid string
	NumMedia      int
	ToCity        string
	FromZip       string
	SmsSid        string
	FromState     string
	SmsStatus     string
	FromCity      string
	Body          string
	FromCountry   string
	To            string
	ToZip         string
	NumSegments   string
	MessageSid    string
	AccountSid    string
	From          string
	ApiVersion    string
}

type SMSServer struct {
	Router    *chi.Mux
	SMSClient *SMSClient
	RootCmd   *cobra.Command
}

// Make the server and it's variables accessible to route
var s *SMSServer

func CreateNewSMSServer(smsClient *SMSClient, cmd *cobra.Command) *SMSServer {
	s := &SMSServer{SMSClient: smsClient, RootCmd: cmd}
	s.Router = chi.NewRouter()
	return s
}

// Given a port start a web server on the port
func InitSMSServer(p string, smsClient *SMSClient, cmd *cobra.Command) {
	s = CreateNewSMSServer(smsClient, cmd)
	s.RegisterHandlers()

	// Launch the server on the specified port
	port := fmt.Sprintf(":%s", p)
	utils.L.L.Fatal(http.ListenAndServe(port, s.Router))
}

func (s *SMSServer) RegisterHandlers() {
	base := "/"

	utils.L.L.Formatter = &logrus.JSONFormatter{
		// disable, as we set our own
		DisableTimestamp: true,
	}

	// Add middleware
	s.Router.Use(middleware.RequestID)
	s.Router.Use(NewStructuredLogger(utils.L.L))
	s.Router.Use(middleware.Heartbeat(base + "/healthcheck"))
	s.Router.Use(middleware.Recoverer)

	// Add routes
	s.Router.Post(base+"receive", ReceiveMessage)
}

func NewTwilioWebhookRequest(r io.ReadCloser) (TwilioWebhookRequestBody, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	if err != nil {
		return TwilioWebhookRequestBody{}, err
	}

	fieldMap, err := url.ParseQuery(buf.String())
	if err != nil {
		return TwilioWebhookRequestBody{}, err
	}

	var incomingMessage TwilioWebhookRequestBody

	flattenedFieldMap := make(map[string]string)
	for k, v := range fieldMap {
		flattenedFieldMap[k] = v[0]
	}

	mapstructure.Decode(flattenedFieldMap, &incomingMessage)
	return incomingMessage, nil
}

// A webhook for processing messages
func ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	incomingMessage, err := NewTwilioWebhookRequest(r.Body)
	if err != nil {
		utils.L.L.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("%+v", incomingMessage)

	destType := Normal
	if utils.IsSat(incomingMessage.From) {
		destType = Satellite
	}

	destination := &SMSDestination{PhoneNumber: incomingMessage.From, Type: destType}

	response := router.RouteStringAsCommand(
		incomingMessage.Body,
		s.RootCmd,
	)

	messages := s.SMSClient.serializeOutbound(response, destination)
	twimlResult := s.SMSClient.parseToTwiML(messages)

	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(twimlResult))
}
