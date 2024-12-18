package service

import (
	"errors"
	"YALP/internal/util"
)

type AuthService interface {
	ValidateJWT(tokenStr string) (int64, error)
}

type authService struct {
	jwtSecret string
}

func NewAuthService(secret string) AuthService {
	return &authService{jwtSecret: secret}
}

// ValidateJWT validates the JWT and returns the user ID.
func (a *authService) ValidateJWT(tokenStr string) (int64, error) {
	userID, err := util.ValidateJWTWithSecret(a.jwtSecret, tokenStr)
	if err != nil {
		return 0, errors.New("token validation failed")
	}
	return userID, nil
}
