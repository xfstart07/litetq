package pool

import "fmt"

var defaultpool = New(32)

// 协程池
type Pool interface {
	Go(func())
}

var _ Pool = (*pool)(nil)

type pool struct {
	queue    chan func()   // 任务队列
	workPool chan struct{} // 带缓存的通道，限制协程创建数量
}

func New(size int) *pool {
	return &pool{
		queue:    make(chan func(), 128),
		workPool: make(chan struct{}, size),
	}
}

func Go(task func()) {
	defaultpool.Go(task)
}

func (p *pool) Go(task func()) {
	// 添加任务到队列中
	select {
	case p.queue <- task:
	}

	select {
	case p.workPool <- struct{}{}:
		go p.worker()
	default:
		// 不需要创建新的 goroutine
		fmt.Println("worker return")
	}

	return
}

func (p *pool) worker() {
	defer func() {
		<-p.workPool
	}()
	for {
		select {
		case task := <-p.queue: // 从队列中获取任务
			task()
		default:
			// No task return
			return
		}
	}
}
