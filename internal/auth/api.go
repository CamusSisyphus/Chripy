package auth

import (
	"errors"
	"net/http"
)

func GetAPIKey(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", errors.New("No authorization found!")
	}
	apiKey := authorization[len("ApiKey "):]
	return apiKey, nil

}
