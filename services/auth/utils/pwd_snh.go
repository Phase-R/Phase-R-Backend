package tools

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
)

func PwdSaltAndHash(password string) (hashedPassword string, err error) {
	hashedPassword, err := hashy(password)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

func hashy(password string) (string, error) {
	salt, err := salted(16)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	b64hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("%s.%s", b64salt, b64hash), nil
}

func salted(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}
