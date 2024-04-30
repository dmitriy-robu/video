package event

import "go-fitness/internal/api/types"

type NotificationEvent struct {
	notification types.Notification
}

func NewNotificationEvent(
	notification types.Notification,
) *NotificationEvent {
	return &NotificationEvent{
		notification: notification,
	}
}

func (e *NotificationEvent) Channel() string {
	return "notification"
}

func (e *NotificationEvent) EventType() string {
	return "new-notification"
}

func (e *NotificationEvent) Data() map[string]interface{} {
	return map[string]interface{}{
		"status":  "",
		"message": "",
	}
}
