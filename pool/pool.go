package pool

var defaultpool = New(32)

// 协程池
type Pool interface {
	Go(func())
	Stop()
}

var _ Pool = (*pool)(nil)

type pool struct {
	queue  chan func()   // 任务队列
	worker chan struct{} // 带缓存的通道，限制协程创建数量
	stop   chan struct{}
}

func New(size int) *pool {
	return &pool{
		queue:  make(chan func(), 128),
		worker: make(chan struct{}, size),
		stop:   make(chan struct{}),
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
	case p.worker <- struct{}{}:
		go p.work()
	default:
		// 不需要创建新的 goroutine
	}

	return
}

func (p *pool) work() {
	defer func() {
		// 工作者退出工作
		<-p.worker
	}()
	for {
		select {
		case task := <-p.queue: // 从队列中获取任务
			task()
		case <-p.stop: // 退出工作
			return
		}
	}
}

func (p *pool) Stop() {
	p.stop <- struct{}{}
}
