package service_test

import (
	"testing"

	"github.com/lwinmgmg/outbox/service"
)

func TestGetMySqlDb(t *testing.T) {
	_, err := service.GetMySqlDb("172.30.1.137", 3301, "lwinmgmg", "letmein", "outbox")
	if err != nil {
		t.Errorf("Error on connecting db : %v\n", err)
	}
}
