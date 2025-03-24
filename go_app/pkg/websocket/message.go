package websocket

import (
	"encoding/json"
	"time"
)

// 消息类型常量
const (
	MessageTypeNotification = "notification"
	MessageTypeActivity    = "activity"
	MessageTypeSystem      = "system"
)

// Message WebSocket消息结构
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	Time int64       `json:"time"`
}

// 消息处理相关方法
func NewMessage(msgType string, data interface{}) *Message {
	return &Message{
		Type: msgType,
		Data: data,
		Time: time.Now().Unix(),
	}
}

func (msg *Message) ToJSON() ([]byte, error) {
	return json.Marshal(msg)
}
