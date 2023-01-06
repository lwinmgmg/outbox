package service_test

import (
	"testing"

	"github.com/lwinmgmg/outbox/service"
)

func TestGetPgDb(t *testing.T) {
	_, err := service.GetPgDb("172.30.1.137", 5432, "lwinmgmg", "letmein", "outbox")
	if err != nil {
		t.Errorf("Error on connecting db : %v\n", err)
	}
}
