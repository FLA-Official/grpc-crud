package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Payload struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func CreateJWT(secret string, data Payload) (string, error) {
	//setting Header (first part of JWT)
	header := Header{
		Alg: "HS256", //may change depending requirment
		Typ: "JWT",   //fixed
	}
	//converting Header to []byte as encoding to Byte64 only accepts []byte
	byteArrHeader, err := json.Marshal(header)
	//error handle
	if err != nil {
		return "", err
	}
	//converting header from []byte to Byte64 string
	headerB64 := Base64UrlEncode(byteArrHeader)

	//converting Payload to []byte
	byteArrData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	//converting Payload from []byte to Byte64 string
	payloadB64 := Base64UrlEncode(byteArrData)
	//turning haeder and payload into single string
	message := headerB64 + "." + payloadB64
	//converting secret from string to byte array
	byteArrSecret := []byte(secret)
	//converting message from string to byte array
	byteArrMessage := []byte(message)
	// using hmac which will take sha256.New (creating new hash by computing SHA256 checksum) adn secret as []byte format
	h := hmac.New(sha256.New, byteArrSecret)
	//converting massage into hash
	h.Write(byteArrMessage)

	signature := h.Sum(nil)
	signatureB64 := Base64UrlEncode(signature)

	jwt := headerB64 + "." + payloadB64 + "." + signatureB64

	return jwt, nil

}

func Base64UrlEncode(data []byte) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

// Base64UrlDecode decodes a base64 URL-encoded string
func Base64UrlDecode(s string) ([]byte, error) {
	// Base64 URL encoding sometimes omits padding, so we need to add it back
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}

	// decode using URLEncoding (uses - and _ instead of + and /)
	return base64.URLEncoding.DecodeString(s)
}
