package entity

import "time"

type VerificationTokenEntity struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Token     string     `json:"token"`
	TokenType string     `json:"token_type"`
	ExpiresAt time.Time  `json:"expires_at"`
	User      UserEntity `json:"user"`
}
