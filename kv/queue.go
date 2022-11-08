package kv

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/xfstart07/litetq"
	"github.com/xfstart07/litetq/kvclient"
)

var _ litetq.QueueClient = new(queue)

type queue struct {
	maxRetryCount int
	db            kvclient.KVClient
	msgch         chan *litetq.Message
	handler       map[litetq.MessageType]func(*litetq.Message) error
	done          chan struct{}
}

func New(kvcli kvclient.KVClient, maxRetryCount int) *queue {
	q := &queue{
		maxRetryCount: maxRetryCount,
		db:            kvcli,
		msgch:         make(chan *litetq.Message, 100),
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
	var err error

	log.Debug().Msgf("dispatch message msgId:%s type:%v", msg.MsgId, msg.Type)
	start := time.Now()
	defer func() {
		log.Info().Msgf("dispatch type(%v) %s done, cost(%d)us.", msg.Type, msg.MsgId, time.Since(start).Microseconds())
	}()

	if msg.Retry <= q.maxRetryCount {
		fn, ok := q.handler[msg.Type]
		if !ok {
			log.Error().Msgf("message type(%s) not subscribe handler", msg.Type)
			return
		}
		err = fn(msg)
	} else {
		log.Error().Msgf("drop message, retry:%d msgId:%s", msg.Retry, msg.MsgId)
		if err = q.db.Set(msg.Table()+msg.MsgId, msg); err != nil {
			log.Error().Msgf("del message db data err: %+v, key:%s", err, msg.MsgId)
			return
		}
	}

	if err == nil { // 事件上报成功，删除缓存
		log.Info().Msgf("message type:%s msgId:%s, run success, to cleanup", msg.Type, msg.MsgId)
		err = q.db.Del(msg.Table() + msg.MsgId)
		if err != nil {
			log.Error().Msgf("del message db data err: %+v, key:%s", err, msg.MsgId)
		}
	} else { // 事件上报失败，更新失败次数，过会再重试
		log.Warn().Msgf("message to retry, type:%s key:%s message:%s, error: %v", msg.Type, msg.MsgId, string(msg.Data), err)
		msg.Retry++
		if err = q.db.Set(msg.Table()+msg.MsgId, msg); err != nil {
			log.Error().Msgf("message db error:%s key:%s", err, msg.MsgId)
		}

		waitTime := msg.Retry * 10 // 秒
		q.publish(msg, waitTime)
	}
}

// Close implements litetq.QueueClient
func (q *queue) Close() {
	close(q.done)
}

// LoadAll implements litetq.QueueClient
func (q *queue) LoadAll() error {
	var m litetq.Message
	list, err := q.db.All(m.Table())
	if err != nil {
		return err
	}

	for _, msg := range list {
		q.publish(msg, 0)
	}
	return nil
}

// Publish implements litetq.QueueClient
func (q *queue) Publish(m *litetq.Message, waitSeconds int) error {
	if err := q.db.Set(m.Table()+m.MsgId, m); err != nil {
		return err
	}

	return q.publish(m, waitSeconds)
}

func (q *queue) publish(m *litetq.Message, waitSeconds int) error {
	if waitSeconds > 0 {
		go time.AfterFunc(time.Duration(waitSeconds)*time.Second, func() {
			q.msgch <- m
		})
	}
	q.msgch <- m
	return nil
}

// Subscribe implements litetq.QueueClient
func (q *queue) Subscribe(t litetq.MessageType, fn func(*litetq.Message) error) {
	q.handler[t] = fn
}
