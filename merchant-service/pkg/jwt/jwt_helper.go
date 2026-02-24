package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"micro-warehouse/merchant-service/configs"
	"strings"
	"time"
)

type JWTClaims struct {
	UserID    int    `json:"user_id"`
	Email     string `json:"email"`
	Roles     string `json:"roles"`
	ExpiredAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
}

func GenerateInternalToken(config configs.Config) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(config.App.JwtDuration) * time.Hour)

	claims := JWTClaims{
		UserID:    0,
		Email:     "system@warehouse.internal",
		Roles:     "system",
		ExpiredAt: expiresAt.Unix(),
		IssuedAt:  now.Unix(),
		Issuer:    config.App.JwtIssuer,
		Subject:   "internal-service-communication",
	}

	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)

	signatureInput := headerEncoded + "." + payloadEncoded

	h := hmac.New(sha256.New, []byte(config.App.JwtSecretKey))
	h.Write([]byte(signatureInput))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	token := signatureInput + "." + signature

	return token, nil
}

func ValidateToken(token string, config configs.Config) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	headerEncoded := parts[0]
	payloadEncoded := parts[1]
	signatureEncoded := parts[2]

	signatureInput := headerEncoded + "." + payloadEncoded

	h := hmac.New(sha256.New, []byte(config.App.JwtSecretKey))
	h.Write([]byte(signatureInput))
	expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	if signatureEncoded != expectedSignature {
		return nil, fmt.Errorf("invalid signature")
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var claims JWTClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if time.Now().Unix() > claims.ExpiredAt {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}
