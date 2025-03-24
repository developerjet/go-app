package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID     uint
	Socket *websocket.Conn
	Send   chan []byte
}

type Manager struct {
	Clients       map[uint]*Client
	Register      chan *Client
	Unregister    chan *Client
	Broadcast     chan []byte
	Subscriptions map[string]map[uint]bool
	mutex         sync.RWMutex
}

// 添加全局管理器
var GlobalManager *Manager

// 保留这个完整实现的 Subscribe 方法
func (m *Manager) Subscribe(userID uint, topic string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.Subscriptions[topic] == nil {
		m.Subscriptions[topic] = make(map[uint]bool)
	}
	m.Subscriptions[topic][userID] = true
}

// 保留这个完整实现的 Unsubscribe 方法
func (m *Manager) Unsubscribe(userID uint, topic string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if subs, exists := m.Subscriptions[topic]; exists {
		delete(subs, userID)
	}
}

// NewManager 创建一个新的 WebSocket 管理器
func NewManager() *Manager {
	return &Manager{
		Clients:       make(map[uint]*Client),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Broadcast:     make(chan []byte),
		Subscriptions: make(map[string]map[uint]bool),
	}
}

// Start 启动 WebSocket 管理器
func (m *Manager) Start() {
	for {
		select {
		case client := <-m.Register:
			m.mutex.Lock()
			m.Clients[client.ID] = client
			m.mutex.Unlock()
			
			// 发送连接成功消息
			m.SendFormattedMessage(client.ID, MessageTypeSystem, map[string]interface{}{
				"message": "WebSocket 连接成功",
				"userId": client.ID,
			})

		case client := <-m.Unregister:
			if _, ok := m.Clients[client.ID]; ok {
				m.mutex.Lock()
				delete(m.Clients, client.ID)
				close(client.Send)
				m.mutex.Unlock()
			}

		case message := <-m.Broadcast:
			m.mutex.RLock()
			for _, client := range m.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(m.Clients, client.ID)
				}
			}
			m.mutex.RUnlock()
		}
	}
}

func (m *Manager) SendFormattedMessage(userID uint, msgType string, data interface{}) error {
	msg := NewMessage(msgType, data)
	jsonData, err := msg.ToJSON()
	if err != nil {
		return err
	}
	m.SendToUser(userID, jsonData)
	return nil
}

func (m *Manager) SendToUser(userID uint, jsonData []byte) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if client, ok := m.Clients[userID]; ok {
		select {
		case client.Send <- jsonData:
			// 消息已发送到客户端的缓冲通道
		default:
			// 客户端缓冲区已满，关闭连接并清理资源
			close(client.Send)
			delete(m.Clients, userID)
		}
	}
}

func (m *Manager) BroadcastToTopic(topic string, msgType string, data interface{}) error {
	msg := NewMessage(msgType, data)
	jsonData, err := msg.ToJSON()
	if err != nil {
		return err
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if subs, exists := m.Subscriptions[topic]; exists {
		for userID := range subs {
			if client, ok := m.Clients[userID]; ok {
				select {
				case client.Send <- jsonData:
					// 消息已发送
				default:
					// 客户端缓冲区已满
					close(client.Send)
					delete(m.Clients, userID)
					delete(subs, userID)
				}
			} else {
				delete(subs, userID)
			}
		}
	}
	return nil
}
