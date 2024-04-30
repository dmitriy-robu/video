package enum

type NotificationStatus int

const (
	NotificationStatusError NotificationStatus = iota
	NotificationStatusWarning
	NotificationStatusSuccess
	NotificationStatusMessage
)

func (e NotificationStatus) String() string {
	switch e {
	case NotificationStatusError:
		return "error"
	case NotificationStatusWarning:
		return "warning"
	case NotificationStatusSuccess:
		return "success"
	case NotificationStatusMessage:
		return "message"
	}
	return ""
}
