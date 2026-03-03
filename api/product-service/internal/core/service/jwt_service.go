package service

import (
	"product-service/config"

	"github.com/golang-jwt/jwt/v5"
)

type JwtServiceInterface interface {
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey string
}

// ValidateToken implements JwtServiceInterface.
func (j *jwtService) ValidateToken(encodeToken string) (*jwt.Token, error) {
	return jwt.Parse(encodeToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(j.secretKey), nil
	})
}

func NewJwtService(cfg *config.Config) JwtServiceInterface {
	return &jwtService{
		secretKey: cfg.App.JwtSecretKey,
	}
}
