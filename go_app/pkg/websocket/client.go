package websocket

import (
    "encoding/json"
    "log"
    "time"

    "github.com/gorilla/websocket"
)

const (
    writeWait      = 10 * time.Second
    pongWait       = 60 * time.Second
    pingPeriod     = (pongWait * 9) / 10
    maxMessageSize = 512
)

// ReadPump 处理从客户端读取消息
func (c *Client) ReadPump() {
    defer func() {
        GlobalManager.Unregister <- c  // 添加这行，确保客户端断开时注销
        c.Socket.Close()
    }()

    c.Socket.SetReadLimit(maxMessageSize)
    c.Socket.SetReadDeadline(time.Now().Add(pongWait))
    c.Socket.SetPongHandler(func(string) error {
        c.Socket.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, message, err := c.Socket.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("读取错误: %v", err)
            }
            break
        }

        // 处理接收到的消息
        var msg struct {
            Action string `json:"action"`
            Topic  string `json:"topic"`
        }
        if err := json.Unmarshal(message, &msg); err != nil {
            log.Printf("消息解析错误: %v", err)
            continue
        }

        // 处理订阅/取消订阅操作
        switch msg.Action {
        case "subscribe":
            GlobalManager.Subscribe(c.ID, msg.Topic)
        case "unsubscribe":
            GlobalManager.Unsubscribe(c.ID, msg.Topic)
        }
    }
}

// WritePump 处理向客户端发送消息
func (c *Client) WritePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.Socket.Close()
    }()

    for {
        select {
        case message, ok := <-c.Send:
            c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := c.Socket.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            if err := w.Close(); err != nil {
                return
            }
        case <-ticker.C:
            c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}