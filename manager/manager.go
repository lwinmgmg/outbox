package manager

import (
	"fmt"
	"time"

	"github.com/lwinmgmg/outbox/config"
	"github.com/lwinmgmg/outbox/helper"
	"github.com/lwinmgmg/outbox/logging"
	"github.com/lwinmgmg/outbox/model"
	"github.com/lwinmgmg/outbox/service"

	"github.com/segmentio/kafka-go"
)

var (
	_logger  = logging.GetLogger()
	settings *config.Config
)

func init() {
	var err error
	settings, err = config.GetConfig("settings.yaml")
	if err != nil {
		panic(fmt.Sprintf("Error on getting settings : %v", err))
	}
}

func RunMgr(conf config.Topic, closeChan chan<- string, isDone *bool) {

	// inform to close channel
	defer func() {
		closeChan <- conf.Name
	}()

	// connecting db
	db, err := service.GetDb(conf.Driver, conf.DbHost, conf.DbPort, conf.DbUser, conf.DbPassword, conf.DbName)
	if err != nil {
		_logger.Error("Error on connecting database in Topic [%v] process\n", conf.Name)
		return
	}

	// get actual db connection to close
	sqlDB, err := db.DB()
	if err != nil {
		_logger.Error("Error on connecting database in Topic [%v] process\n", conf.Name)
		return
	}

	// close the db connection
	defer sqlDB.Close()

	// getting kafka producer
	producer := service.NewProducer(conf.Brokers, conf.Name, conf.KafkaUser, conf.KafkaPassword, time.Millisecond*3000)
	// close kafka connection on exit
	defer producer.Writer.Close()

	for !*isDone {
		func() {
			// start the db transaction
			tx := db.Begin()

			// rollback on panic
			defer func() {
				if r := recover(); r != nil {
					_logger.Error("Getting panic on transaction and rolled back the transaction : %v", r)
					tx.Rollback()
				}
			}()
			if tx.Error != nil {
				return
			}
			ids := make([]int, 0, 10)
			if err := tx.Model(model.Outbox{}).Select("id").Order(conf.Order).Where("name = ? AND is_produced = ?", conf.Name, false).Limit(conf.ProduceCount).Find(&ids).Error; err != nil {
				_logger.Error("Error on executing query : %v", err)
				tx.Rollback()
				return
			}
			if len(ids) == 0 {
				tx.Rollback()
				return
			}
			// outbox produce
			model.GetByIDForUpdate(tx, ids, func(dataList []model.Outbox) error {
				mesgList := make([]kafka.Message, 0, 10)
				for _, data := range dataList {
					mesg, err := data.PrepareMessage()
					if err != nil {
						_logger.Error("Error on preparing message : ", err)
						_logger.Error("#Skipped - Table ID: %v, Topic Name: %v, Key: %v, Header: %v, Value: %v", data.ID, data.Name, data.Key, data.Headers, data.Value)
						// mute the err not to block the next process
						continue
					}
					mesgList = append(mesgList, mesg)
				}
				if err := helper.Retry(settings.Server.RetryCount, producer.Produce, mesgList...); err != nil {
					return err
				}
				return nil
			})
		}()
		time.Sleep(time.Millisecond * time.Duration(conf.ProduceIntervalMs))
	}
}
