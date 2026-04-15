package user

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func comparePassword(storedHash, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(plain))
	return err == nil
}

func hashPassword(plain string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func randomPassword() string {
	b := make([]byte, 32)
	rand.Read(b)

	raw := base64.URLEncoding.EncodeToString(b)

	// hash string random tersebut
	hashed, _ := hashPassword(raw)

	return hashed
}
