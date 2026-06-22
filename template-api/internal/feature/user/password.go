package user

import "golang.org/x/crypto/bcrypt"

// hashPassword returns a bcrypt hash of the plaintext password.
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// checkPassword reports whether the plaintext matches the bcrypt hash.
func checkPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
