package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Claims struct {
	Exp int64 `json:"exp"`
}

func IsJWTExpired(token string) (bool, error) {
	// Split the token into its parts (header, payload, signature)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, fmt.Errorf("invalid token format")
	}

	// Base64 decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false, fmt.Errorf("error decoding payload: %v", err)
	}

	// Parse the payload into a Claims struct
	var claims Claims
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return false, fmt.Errorf("error unmarshalling claims: %v", err)
	}

	// Get the current time and compare it with the exp field
	now := time.Now().Unix()
	if claims.Exp < now {
		return true, nil // Token is expired
	}

	return false, nil // Token is not expired
}
