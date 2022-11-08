package kv

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/xfstart07/litetq"
	"github.com/xfstart07/litetq/kvclient/leveldb"
)

func TestMain(m *testing.M) {
	m.Run()

	os.RemoveAll("./data")
}

const taskA litetq.MessageType = "taskA"

func taskAFn(msg *litetq.Message) error {
	log.Info().Msgf("taskAAA: %+v", msg)
	return nil
}

const taskB litetq.MessageType = "taskB"

func taskBFn(msg *litetq.Message) error {
	log.Info().Msgf("taskBBBL: %+v", msg)
	return nil
}

func TestQueue(t *testing.T) {
	db, err := leveldb.New("./data")
	if err != nil {
		t.Fatal(err)
	}

	queue := New(db, 10)
	queue.Subscribe(taskA, taskAFn)
	queue.Subscribe(taskB, taskBFn)

	for i := 0; i < 10000; i++ {
		if i%3 == 0 {
			msg := litetq.NewMessage(taskA, "test")
			queue.Publish(msg, 0)
		} else {

			msg := litetq.NewMessage(taskB, "testBBB")
			queue.Publish(msg, 0)
		}
	}
	time.Sleep(2 * time.Second)
}
