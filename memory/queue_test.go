package memory

import (
	"math/rand"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/xfstart07/litetq"
)

const taskA litetq.MessageType = "taskA"

func taskAFn(msg *litetq.Message) error {
	log.Info().Msgf("taskAAA: %+v", msg)
	time.Sleep(10 * time.Microsecond)
	return nil
}

const taskB litetq.MessageType = "taskB"

func taskBFn(msg *litetq.Message) error {
	log.Info().Msgf("taskBBBL: %+v", msg)
	return nil
}

func BenchmarkNew(b *testing.B) {
	queue := New(30)
	queue.Subscribe(taskA, taskAFn)
	queue.Subscribe(taskB, taskBFn)

	for i := 0; i < b.N; i++ {
		n := rand.Intn(2)
		if n == 0 {
			msg := litetq.NewMessage(taskA, "test")
			queue.Publish(msg, 0)
		}
	}

}
