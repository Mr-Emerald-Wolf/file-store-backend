package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType int

const (
	ACCESS_TOKEN = iota
	REFRESH_TOKEN
)

func CreateToken(email string, tokenType TokenType) (string, error) {

	secret := []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	expiry := time.Now().Add(time.Minute * 15).Unix()

	if tokenType != ACCESS_TOKEN {
		secret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))
		expiry = time.Now().Add(time.Hour).Unix()
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,             // Subject (user identifier)
		"iss": "file-share-app",  // Issuer
		"aud": "user",            // Audience (user role)
		"exp": expiry,            // Expiration time
		"iat": time.Now().Unix(), // Issued at
	})

	tokenString, err := claims.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil

}
