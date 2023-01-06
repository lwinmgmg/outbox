package logging

import (
	"github.com/lwinmgmg/logger"
)

var logging *logger.Logging = nil

func GetLogger() *logger.Logging {
	if logging == nil {
		logging = logger.DefaultLogging(logger.INFO)
	}
	return logging
}
