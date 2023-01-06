package logging_test

import (
	"testing"

	"github.com/lwinmgmg/outbox/logging"
)

func TestGetLogger(t *testing.T) {
	_logger := logging.GetLogger()
	_logger.Info("INFO")
	_logger.Warning("WARNING")
	_logger.Error("ERROR")
}
