package entity

type WsRedisEntity struct {
	ID         int64  `json:"id"`
	ReceiverID int64  `json:"receiver_id"`
	Type       string `json:"type"`
	Subject    string `json:"subject"`
	Message    string `json:"message"`
	SentAt     string `json:"sent_at"`
}
