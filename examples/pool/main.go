package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	_ "net/http/pprof"

	"github.com/xfstart07/litetq/pool"
)

func main() {
	go func() {
		log.Fatal(http.ListenAndServe(":6060", nil))
	}()

	start := time.Now().UnixMilli()

	wg := sync.WaitGroup{}
	for i := 0; i < 100000; i++ {
		wg.Add(1)

		pool.Go(func() {
			fmt.Println("run...........")
			wg.Done()
		})
	}
	log.Println("wait..")
	wg.Wait()
	time.Sleep(time.Second)
	fmt.Println("time spent:", time.Now().UnixMilli()-start, "ms")
}
