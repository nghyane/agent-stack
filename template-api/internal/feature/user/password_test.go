package user

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	const pw = "password123"

	hash, err := hashPassword(pw)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if hash == pw {
		t.Fatal("hash must not equal plaintext")
	}
	if !checkPassword(hash, pw) {
		t.Error("correct password should match")
	}
	if checkPassword(hash, "wrong") {
		t.Error("wrong password should not match")
	}
}
