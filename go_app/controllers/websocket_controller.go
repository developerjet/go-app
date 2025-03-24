package controllers

import (
    "go_app/models"
    "go_app/pkg/errcode"
    "go_app/pkg/websocket"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    gorillaws "github.com/gorilla/websocket"
)

type WebSocketController struct {
    manager *websocket.Manager
}

func NewWebSocketController(manager *websocket.Manager) *WebSocketController {
    return &WebSocketController{
        manager: manager,
    }
}

func (wc *WebSocketController) HandleConnection(c *gin.Context) {
    userID, exists := c.Get("userId")
    if !exists {
        log.Printf("用户认证失败")
        c.JSON(http.StatusUnauthorized, models.NewError(errcode.Unauthorized))
        return
    }

    // 检查用户是否已经连接
    if _, ok := wc.manager.Clients[userID.(uint)]; ok {
        log.Printf("用户 %d 已经连接", userID)
        c.JSON(http.StatusBadRequest, models.NewError(errcode.InvalidRequest))
        return
    }

    upgrader := gorillaws.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("WebSocket 升级失败: %v", err)
        return
    }

    log.Printf("WebSocket 连接成功建立，用户ID: %v", userID)

    client := &websocket.Client{
        ID:     userID.(uint),
        Socket: conn,
        Send:   make(chan []byte, 256),
    }

    wc.manager.Register <- client

    go client.ReadPump()
    go client.WritePump()
}