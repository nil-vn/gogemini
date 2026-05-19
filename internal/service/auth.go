package service

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

func CheckWerkzeugPasswordHash(encoded, password string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 3 { return false }
	method, salt, expected := parts[0], parts[1], parts[2]
	m := strings.Split(method, ":")
	if len(m) < 2 || m[0] != "pbkdf2" || m[1] != "sha256" { return false }
	iterations := 260000
	if len(m) >= 3 { if n, err := strconv.Atoi(m[2]); err == nil { iterations = n } }
	derived := pbkdf2.Key([]byte(password), []byte(salt), iterations, sha256.Size, sha256.New)
	actual := base64.StdEncoding.EncodeToString(derived)
	return subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) == 1
}

func BuildSessionToken(userID int64, secret string) string {
	return fmt.Sprintf("u:%d", userID)
}
