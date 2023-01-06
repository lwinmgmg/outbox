package model

import (
	"encoding/json"
	"errors"

	"github.com/lwinmgmg/outbox/logging"
	"github.com/lwinmgmg/outbox/service"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	_logger = logging.GetLogger()
)

type Outbox struct {
	ID         int
	Name       string
	Key        string
	Value      string
	Headers    string
	IsProduced bool
}

func (outbox *Outbox) TableName() string {
	return "outbox"
}

func (outbox Outbox) PrepareMessage() (kafka.Message, error) {
	data := make(map[string]any, 10)
	if outbox.Headers != "" {
		if err := json.Unmarshal([]byte(outbox.Headers), &data); err != nil {
			return kafka.Message{}, err
		}
		if len(data) == 0 {
			return kafka.Message{}, errors.New("error on preparing header")
		}
	}
	return kafka.Message{
		Key:     []byte(outbox.Key),
		Value:   []byte(outbox.Value),
		Headers: service.MapToHeaders(data),
	}, nil
}

func GetByIDForUpdate(tx *gorm.DB, ids []int, callBack func([]Outbox) error) {
	var outboxList []Outbox
	if err := tx.Clauses(clause.Locking{
		Strength: "UPDATE",
		Options:  "NOWAIT",
	}).Where("id in ?", ids).Find(&outboxList).Error; err != nil {
		_logger.Error("Error on locking the table : ", err)
		tx.Rollback()
		return
	}
	if err := callBack(outboxList); err != nil {
		_logger.Error("Error on running producing : ", err)
		tx.Rollback()
		return
	}
	if err := tx.Model(Outbox{}).Where("id IN ?", ids).Updates(Outbox{IsProduced: true}).Error; err != nil {
		_logger.Error("Error on updating is_produced status : ", err)
		tx.Rollback()
		return
	}
	_logger.Info("Commited")
	tx.Commit()
}
