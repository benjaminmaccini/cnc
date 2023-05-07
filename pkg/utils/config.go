package utils

type Config struct {
	BIN_SIZE           int // Char limit for satellite phones
	OPENAI_KEY         string
	SMS_SERVER_PORT    string
	TWILIO_ACCOUNT_SID string
	TWILIO_KEY         string
}

// Create the globally available config instance
var C *Config
