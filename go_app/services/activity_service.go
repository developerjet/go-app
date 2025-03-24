package services

import (
    "go_app/models"
    "go_app/pkg/websocket"
)

type ActivityService struct {
    wsManager *websocket.Manager
}

func NewActivityService(wsManager *websocket.Manager) *ActivityService {
    return &ActivityService{
        wsManager: wsManager,
    }
}

func (s *ActivityService) PushActivity(activity *models.Activity) error {
    return s.wsManager.BroadcastToTopic("activity", websocket.MessageTypeActivity, activity)
}

func (s *ActivityService) PushNotification(userID uint, notification *models.Notification) error {
    return s.wsManager.SendFormattedMessage(userID, websocket.MessageTypeNotification, notification)
}