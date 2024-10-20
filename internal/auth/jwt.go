package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
)

type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "chirpy-access"
)

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	var secretKey = []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    string(TokenTypeAccess),
		Subject:   userID.String(),
	})
	ss, err := token.SignedString(secretKey)
	return ss, err

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Printf("Error parsing JWT: %v", err)
		return uuid.Nil, fmt.Errorf("Error parsing JWT: %w", err)
	}

	if !token.Valid {
		log.Printf("Invalid token")
		return uuid.Nil, fmt.Errorf("Invalid token")
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return uuid.Nil, fmt.Errorf("Couldn't get expiraton time: %w", err)
	}
	if exp.Add(-time.Minute).Before(time.Now()) {
		return uuid.Nil, fmt.Errorf("JWT expired: %w", err)
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	BT := headers.Get("Authorization")
	if BT == "" {
		return "", errors.New("No Bearer Token")
	}
	tokenString := BT[len("Bearer "):]
	return tokenString, nil
}
