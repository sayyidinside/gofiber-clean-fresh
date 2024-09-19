package worker

import (
	"log"
	"time"

	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type LogData struct {
	Timestamp time.Time
	Message   string
	LogType   string
}

var LogChannel = make(chan LogData, 100)

func StartLogWorker() {
	go func() {
		for {
			select {
			case logSys := <-helpers.LogSysChannel:
				log.Println(logSys)
				helpers.GenerateLogSystem(logSys)
			case logAPI := <-helpers.LogAPIChannel:
				log.Println(logAPI)
				helpers.GenerateLogAPI(logAPI)
			}
		}
	}()
}
