package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lwinmgmg/outbox/config"
	"github.com/lwinmgmg/outbox/logging"
	"github.com/lwinmgmg/outbox/manager"

	"github.com/lwinmgmg/logger"
)

var (
	settings *config.Config  = nil
	_logger  *logger.Logging = logging.GetLogger()
)

func init() {
	var err error
	_logger.Info("Started")
	settings, err = config.GetConfig("settings.yaml")
	if err != nil {
		panic(fmt.Sprintf("Error on getting settings : %v", err))
	}
}

func main() {
	// Declare variables
	var isDone bool = false
	processChan := make(chan string)
	doneChan := make(chan struct{}, 1)

	// Close Signal
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		isDone = true
		_logger.Info("Close by user")
	}()

	topicMap := make(map[string]config.Topic, 10)
	for _, topic := range settings.Server.Topics {
		topicMap[topic.Name] = topic
		go manager.RunMgr(topic, processChan, &isDone)
	}

	go func() {
		for {
			if topic, ok := <-processChan; ok {
				// respawn the thread when failed after 5 sec
				if isDone {
					delete(topicMap, topic)
					if len(topicMap) > 0 {
						continue
					}
					break
				}
				time.Sleep(time.Millisecond * time.Duration(settings.Server.RespawnTimeMs))
				go manager.RunMgr(topicMap[topic], processChan, &isDone)
			} else {
				break
			}
		}
		close(doneChan)
	}()
	<-doneChan
	_logger.Info("Successfully Closed")
	_logger.Close()
}
