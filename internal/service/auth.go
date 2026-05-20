package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

func BuildWerkzeugPasswordHash(password string) (string, error) {
	const iterations = 260000
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	saltB64 := base64.RawURLEncoding.EncodeToString(salt)
	derived := pbkdf2.Key([]byte(password), []byte(saltB64), iterations, sha256.Size, sha256.New)
	hash := base64.StdEncoding.EncodeToString(derived)
	return fmt.Sprintf("pbkdf2:sha256:%d$%s$%s", iterations, saltB64, hash), nil
}

func CheckWerkzeugPasswordHash(encoded, password string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 3 {
		return false
	}
	method, salt, expected := parts[0], parts[1], parts[2]
	m := strings.Split(method, ":")
	if len(m) < 2 || m[0] != "pbkdf2" || m[1] != "sha256" {
		return false
	}
	iterations := 260000
	if len(m) >= 3 {
		if n, err := strconv.Atoi(m[2]); err == nil {
			iterations = n
		}
	}
	derived := pbkdf2.Key([]byte(password), []byte(salt), iterations, sha256.Size, sha256.New)
	actual := base64.StdEncoding.EncodeToString(derived)
	return subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) == 1
}

func BuildSessionToken(userID int64, secret string) string {
	expiresAt := time.Now().UTC().Add(24 * time.Hour).Unix()
	payload := fmt.Sprintf("%d:%d", userID, expiresAt)
	sig := signPayload(payload, secret)
	return payload + ":" + sig
}

func ValidateSessionToken(token, secret string) (int64, bool) {
	parts := strings.Split(token, ":")
	if len(parts) != 3 {
		return 0, false
	}
	payload := parts[0] + ":" + parts[1]
	if subtle.ConstantTimeCompare([]byte(signPayload(payload, secret)), []byte(parts[2])) != 1 {
		return 0, false
	}
	uid, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || uid <= 0 {
		return 0, false
	}
	expiresAt, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || time.Now().UTC().Unix() > expiresAt {
		return 0, false
	}
	return uid, true
}

func signPayload(payload, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	_, _ = h.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
