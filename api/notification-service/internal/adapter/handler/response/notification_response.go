package response

type NotificationResponseList struct {
	ID                 int64  `json:"id"`
	NotificationTypeID *int64 `json:"notification_type_id"`
	NotificationType   string `json:"notification_type"`
	Subject            string `json:"subject"`
	Message            string `json:"message"`
	Status             string `json:"status"`
	SentAt             string `json:"sent_at"`
}

type NotificationDetailResponse struct {
	ID                 int64  `json:"id"`
	Subject            string `json:"subject"`
	Message            string `json:"message"`
	Status             string `json:"status"`
	SentAt             string `json:"sent_at"`
	ReadAt             string `json:"read_at"`
	NotificationMethod string `json:"notification_method"`
}
