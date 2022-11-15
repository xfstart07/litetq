package pool

import "github.com/rs/zerolog/log"

type Pool interface {
	Run(func())
	Stop()
}

var _ Pool = (*pool)(nil)

type pool struct {
	queue  chan func()
	worker chan struct{}
	stop   chan struct{}
}

func New(size int) *pool {
	return &pool{
		queue:  make(chan func(), 1),
		worker: make(chan struct{}, size),
		stop:   make(chan struct{}),
	}
}

func (p *pool) Run(task func()) {
	select {
	case p.queue <- task:
	case p.worker <- struct{}{}:
		go p.work()
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
			log.Info().Msg("work done.")
			return
		}
	}
}

func (p *pool) Stop() {
	p.stop <- struct{}{}
}
