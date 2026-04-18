package support

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Andras5014/EIBotHub/server/internal/model"
)

type TokenClaims struct {
	UserID uint
	Role   string
}

type TokenManager struct {
	secret string
}

func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{secret: secret}
}

func (m *TokenManager) HashPassword(password string) string {
	sum := sha256.Sum256([]byte(m.secret + ":" + password))
	return hex.EncodeToString(sum[:])
}

func (m *TokenManager) CheckPassword(hashedPassword, password string) bool {
	expected := m.HashPassword(password)
	return subtle.ConstantTimeCompare([]byte(hashedPassword), []byte(expected)) == 1
}

func (m *TokenManager) Issue(user model.User) string {
	expiry := time.Now().Add(7 * 24 * time.Hour).Unix()
	payload := fmt.Sprintf("%d|%s|%d", user.ID, user.Role, expiry)
	signature := m.sign(payload)
	raw := payload + "|" + signature
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func (m *TokenManager) Parse(token string) (TokenClaims, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return TokenClaims{}, NewError(401, "invalid_token", "invalid token")
	}

	parts := strings.Split(string(decoded), "|")
	if len(parts) != 4 {
		return TokenClaims{}, NewError(401, "invalid_token", "invalid token")
	}

	payload := strings.Join(parts[:3], "|")
	if !hmac.Equal([]byte(parts[3]), []byte(m.sign(payload))) {
		return TokenClaims{}, NewError(401, "invalid_token", "invalid token")
	}

	userID64, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return TokenClaims{}, NewError(401, "invalid_token", "invalid token")
	}

	expiry, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || time.Now().Unix() > expiry {
		return TokenClaims{}, NewError(401, "token_expired", "token expired")
	}

	return TokenClaims{
		UserID: uint(userID64),
		Role:   parts[1],
	}, nil
}

func (m *TokenManager) sign(payload string) string {
	mac := hmac.New(sha256.New, []byte(m.secret))
	_, _ = mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func SplitCSV(input string) []string {
	rawParts := strings.Split(input, ",")
	result := make([]string, 0, len(rawParts))
	for _, part := range rawParts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func JoinCSV(values []string) string {
	clean := make([]string, 0, len(values))
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			clean = append(clean, trimmed)
		}
	}
	return strings.Join(clean, ",")
}
