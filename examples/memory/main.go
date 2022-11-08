package main

import (
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/xfstart07/litetq"
	"github.com/xfstart07/litetq/memory"
)

var queue litetq.QueueClient

func main() {
	go func() {
		log.Fatal().Err(http.ListenAndServe("localhost:6060", nil))
	}()

	queue = memory.New(100, 10)

	queue.Subscribe(reportAlert, reportAlertFn)
	queue.Subscribe(sendMail, sendMailFn)

	tick := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-tick.C:
			mock()
		}
	}
}

func mock() {
	type alertType struct {
		Alert   string
		Content string
	}

	for i := 0; i < 500; i++ {
		n := rand.Intn(2)
		if n == 0 {
			msg := litetq.NewMessage(reportAlert, alertType{Alert: "ERROR", Content: "Failed to connection cloud."})
			queue.Publish(msg, 2)
		} else {
			msg := litetq.NewMessage(sendMail, "123456789@qq.com")
			queue.Publish(msg, 0)
		}
	}
}

const reportAlert litetq.MessageType = "reportAlert"

func reportAlertFn(m *litetq.Message) error {
	log.Info().Msgf("reportAlertFn: %v, msgId: %s", m.Data, m.MsgId)
	n := rand.Intn(100)
	if n == 0 {
		return fmt.Errorf("reportAlert error: connection fail")
	}
	time.Sleep(time.Duration(n) * time.Microsecond)
	return nil
}

const sendMail litetq.MessageType = "sendmail"

func sendMailFn(m *litetq.Message) error {
	log.Info().Msgf("send mail: %v, msgId: %s", m.Data, m.MsgId)
	n := rand.Intn(100)
	time.Sleep(time.Duration(n) * time.Microsecond)
	return nil
}
