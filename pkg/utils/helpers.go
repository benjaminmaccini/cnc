package utils

import (
	"strings"

	"github.com/mattn/go-shellwords"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/constraints"
)

// Parse a string as if it were a command
func ParseStringAsCommand(input string) []string {
	args, err := shellwords.Parse(input)
	if err != nil {
		log.Fatal(err)
	}
	return args
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Check if a phone number is satellite or not based on the country code
// This requires the numbers to be E.614 formatted
// https://en.wikipedia.org/wiki/Global_Mobile_Satellite_System
func IsSat(number string) bool {
	if strings.HasPrefix(number, "+881") {
		return true
	}

	return false
}
