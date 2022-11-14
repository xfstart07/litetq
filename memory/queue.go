package memory

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/xfstart07/litetq"
)

var _ litetq.QueueClient = (*queue)(nil)

type queue struct {
	maxRetryCount int
	msgch         chan *litetq.Message
	handler       map[litetq.MessageType]func(*litetq.Message) error
	done          chan struct{}
}

// 内存队列
// 队列容量 10-1000 范围
func New(concurrent int, maxRetryCount int) *queue {
	if concurrent < 10 {
		concurrent = 10
	}
	if concurrent > 1000 {
		concurrent = 1000
	}
	q := &queue{
		maxRetryCount: maxRetryCount,
		msgch:         make(chan *litetq.Message, concurrent),
		handler:       make(map[litetq.MessageType]func(*litetq.Message) error),
		done:          make(chan struct{}),
	}

	go q.run()

	return q
}

func (q *queue) run() {
	for {
		select {
		case msg := <-q.msgch:
			q.dispatch(msg)
		case <-q.done:
			log.Info().Msg("exit queue!")
			return
		}
	}
}

func (q *queue) dispatch(msg *litetq.Message) {
	log.Info().Msgf("dispatch msg: %+v", msg)
	start := time.Now()
	defer func() {
		log.Info().Msgf("dispatch type(%v) %s done, cost(%d)us.", msg.Type, msg.MsgId, time.Since(start).Microseconds())
	}()

	var err error
	if msg.Retry <= q.maxRetryCount {
		fn, ok := q.handler[msg.Type]
		if !ok {
			log.Error().Msgf("dispatch msgType(%v) %s unknown, data: %+v", msg.Type, msg.MsgId, msg.Data)
		}
		err := fn(msg)
		if err != nil {
			log.Error().Msgf("execute fn error: %v", err)
		}
	} else {
		log.Error().Msgf("drop message, retry:%d, msgId:%s", msg.Retry, msg.MsgId)
		return
	}

	if err == nil {
		log.Info().Msgf("message type:%s msgId:%s, run success, to cleanup", msg.Type, msg.MsgId)
		return
	} else { // 事件上报失败，更新失败次数，过会再重试
		log.Warn().Msgf("message to retry, type:%s key:%s message:%s, error: %v", msg.Type, msg.MsgId, string(msg.Data), err)
		msg.Retry++

		waitTime := msg.Retry * 10 // 秒
		q.Publish(msg, waitTime)
	}
}

// Subscribe 订阅消息类型处理函数
func (q *queue) Subscribe(t litetq.MessageType, fn func(*litetq.Message) error) {
	q.handler[t] = fn
}

// Publish 带延迟的消息发布
func (q *queue) Publish(m *litetq.Message, waitSeconds int) error {
	if waitSeconds > 0 {
		go time.AfterFunc(time.Duration(waitSeconds)*time.Second, func() {
			q.msgch <- m
		})
	}
	q.msgch <- m
	return nil
}

// LoadAll 加载持久化消息
func (q *queue) LoadAll() error {
	// 不支持持久化
	return nil
}

// Close 退出消息总线
func (q *queue) Close() {
	close(q.done)
}
