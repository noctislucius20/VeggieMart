package response

type NotificationResponseList struct {
	ID      uint   `json:"id"`
	Subject string `json:"subject"`
	Status  string `json:"status"`
	SentAt  string `json:"sent_at"`
}

type NotificationDetailResponse struct {
	ID               uint   `json:"id"`
	Subject          string `json:"subject"`
	Message          string `json:"message"`
	Status           string `json:"status"`
	SentAt           string `json:"sent_at"`
	ReadAt           string `json:"read_at"`
	NotificationType string `json:"notification_type"`
}
