package model_test

import (
	"testing"

	"github.com/lwinmgmg/outbox/model"
	"github.com/lwinmgmg/outbox/service"

	"gorm.io/gorm"
)

type Table struct {
	gorm.Model
	model.Outbox
}

func (tb *Table) TableName() string {
	return "outbox"
}

func TestCreateTable(t *testing.T) {
	db, err := service.GetPgDb("172.30.1.137", 5432, "lwinmgmg", "letmein", "outbox")
	if err != nil {
		t.Errorf("Error on connection : %v", err)
	}
	db.Migrator().AutoMigrate(&Table{})
}
