package config_test

import (
	"testing"

	"github.com/lwinmgmg/outbox/config"
)

func TestGetConfig(t *testing.T) {
	_, err := config.GetConfig("../settings.yaml.example")
	if err != nil {
		t.Errorf("Error on getting config : %v", err)
	}
}
