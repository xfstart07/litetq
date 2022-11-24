package pool

import (
	"fmt"
	"sync"
	"testing"
)

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
