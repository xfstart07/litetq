package litetq

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const MaxRetrySeconds = 300

type MessageType string

type Message struct {
	MsgId   string      `json:"msgId"`   // 消息ID
	Retry   int         `json:"retry"`   // 重试次数
	Type    MessageType `json:"type"`    // 消息类型
	Created time.Time   `json:"created"` // 事件时间
	Data    string      `json:"data"`    // 消息内容,JSON
}

func (m *Message) Table() string {
	return "tqmsg"
}

func (m *Message) String() string {
	return fmt.Sprintf("msgId:%s, retry:%d, type:%s, created:%s, data:%s",
		m.MsgId, m.Retry, m.Type, m.Created, m.Data)
}

func (m *Message) Unmarshal(v interface{}) error {
	return json.Unmarshal([]byte(m.Data), &v)
}

// NewMessage 创建Message, data会JSON
func NewMessage(t MessageType, data interface{}) *Message {
	b, _ := json.Marshal(data)

	return &Message{
		MsgId:   NewId(),
		Type:    t,
		Created: time.Now(),
		Data:    string(b),
	}
}

func NewId() string {
	return strings.Replace(uuid.NewString(), "-", "", -1)
}
