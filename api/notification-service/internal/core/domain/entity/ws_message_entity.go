package entity

type WsRedisEntity struct {
	ID                 int64  `json:"id"`
	NotificationType   string `json:"notification_type"`
	NotificationTypeID int64  `json:"notification_type_id"`
	NotificationMethod string `json:"notification_method"`
	ReceiverID         int64  `json:"receiver_id"`
	Subject            string `json:"subject"`
	Message            string `json:"message"`
	SentAt             string `json:"sent_at"`
}
