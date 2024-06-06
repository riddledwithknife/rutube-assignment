package testutils

import (
	"github.com/dgrijalva/jwt-go"
	"rutube-assignment/internal/config"
	"time"
)

func GenerateMockToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.GetJWTSecretKey()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
