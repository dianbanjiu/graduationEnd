package common

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("GDxAp6GoyJAZS9jwsNHzUXFq6ePzvKCf")

type Claims struct {
	UserID string `json:"userid"`
	jwt.StandardClaims
}

func ReleaseToken(id string) string {
	claims := &Claims{
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}

func CheckToken(tokenString string) (*jwt.Token, *Claims, error) {
	var claims = &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token, claims, err
}
