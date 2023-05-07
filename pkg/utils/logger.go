package utils

import (
	log "github.com/sirupsen/logrus"
)

type BaseLogger struct {
	L *log.Logger
}

var L *BaseLogger
