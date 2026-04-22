package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"strings"
)

func VerifyJWT(secret string, token string) (*Payload, error) {

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	header := parts[0]
	payload := parts[1]
	signature := parts[2]

	message := header + "." + payload

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	expectedSignature := Base64UrlEncode(h.Sum(nil))

	if expectedSignature != signature {
		return nil, errors.New("invalid signature")
	}

	// decode payload
	payloadBytes, err := Base64UrlDecode(payload)
	if err != nil {
		return nil, err
	}

	var data Payload
	err = json.Unmarshal(payloadBytes, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
