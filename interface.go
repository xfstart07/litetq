package litetq

// 消息总线
type QueueClient interface {
	// Subscribe 订阅消息类型处理函数
	Subscribe(t MessageType, fn func(*Message) error)
	// Publish 带延迟的消息发布
	Publish(m *Message, waitSeconds int) error
	// LoadAll 加载持久化消息
	LoadAll() error
	// Close 退出消息总线
	Close()
}
