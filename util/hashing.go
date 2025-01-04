package util

import "golang.org/x/crypto/bcrypt"

func GenerateHashString(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hashedString := string(hashedBytes[:])
	return hashedString, nil
}

func CompareHashString(hasedString string, s string) error {
	incoming := []byte(s)
	existing := []byte(hasedString)

	return bcrypt.CompareHashAndPassword(existing, incoming)
}
