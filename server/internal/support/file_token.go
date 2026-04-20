package support

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

type FileTokenClaims struct {
	FileID uint
	Inline bool
}

type FileTokenManager struct {
	secret string
}

func NewFileTokenManager(secret string) *FileTokenManager {
	return &FileTokenManager{secret: secret}
}

func (m *FileTokenManager) Issue(fileID uint, inline bool, ttl time.Duration) string {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	inlineFlag := "0"
	if inline {
		inlineFlag = "1"
	}
	expiry := time.Now().Add(ttl).Unix()
	payload := strconv.FormatUint(uint64(fileID), 10) + "|" + inlineFlag + "|" + strconv.FormatInt(expiry, 10)
	signature := m.sign(payload)
	raw := payload + "|" + signature
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func (m *FileTokenManager) Parse(token string) (FileTokenClaims, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return FileTokenClaims{}, NewError(401, "invalid_file_token", "invalid file token")
	}

	parts := strings.Split(string(decoded), "|")
	if len(parts) != 4 {
		return FileTokenClaims{}, NewError(401, "invalid_file_token", "invalid file token")
	}

	payload := strings.Join(parts[:3], "|")
	if !hmac.Equal([]byte(parts[3]), []byte(m.sign(payload))) {
		return FileTokenClaims{}, NewError(401, "invalid_file_token", "invalid file token")
	}

	fileID, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return FileTokenClaims{}, NewError(401, "invalid_file_token", "invalid file token")
	}
	expiry, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || time.Now().Unix() > expiry {
		return FileTokenClaims{}, NewError(401, "file_token_expired", "file token expired")
	}

	return FileTokenClaims{
		FileID: uint(fileID),
		Inline: parts[1] == "1",
	}, nil
}

func (m *FileTokenManager) sign(payload string) string {
	mac := hmac.New(sha256.New, []byte(m.secret))
	_, _ = mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
