package pool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGoPool(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		Go(func() {
			time.Sleep(10 * time.Microsecond)
			wg.Done()
		})
	}
	wg.Wait()
	time.Sleep(time.Second)
}

func BenchmarkGo(b *testing.B) {
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		Go(func() {
			fmt.Println("run")
			wg.Done()
		})
	}
	wg.Wait()
}
