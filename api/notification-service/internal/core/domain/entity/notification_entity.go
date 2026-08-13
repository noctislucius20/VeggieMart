package entity

import "time"

type NotificationEntity struct {
	ID                 int64      `json:"id"`
	NotificationType   string     `json:"notification_type"`
	NotificationTypeID *int64     `json:"notification_type_id"`
	NotificationMethod string     `json:"notification_method"`
	ReceiverID         *int64     `json:"receiver_id"`
	Subject            *string    `json:"subject"`
	Message            string     `json:"message"`
	ReceiverEmail      *string    `json:"receiver_email"`
	SentAt             *time.Time `json:"sent_at"`
	ReadAt             *time.Time `json:"read_at"`
	Status             string     `json:"status"`
}

type NotificationQueryString struct {
	Page      int64
	Limit     int64
	Search    string
	Status    string
	OrderBy   string
	OrderType string
	UserID    int64
	IsRead    bool
}
