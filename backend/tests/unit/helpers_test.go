package unit

import (
	"finara-backend/internal/utils"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if hashedPassword == "" {
		t.Error("Expected hashed password, got empty string")
	}

	if hashedPassword == password {
		t.Error("Hashed password should not equal original password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"
	hashedPassword, _ := utils.HashPassword(password)

	// Test correct password
	if !utils.CheckPassword(password, hashedPassword) {
		t.Error("Expected password to match")
	}

	// Test incorrect password
	if utils.CheckPassword("wrongpassword", hashedPassword) {
		t.Error("Expected password to not match")
	}
}

func TestCheckPasswordHash_WithInvalidHash(t *testing.T) {
	password := "testpassword123"
	invalidHash := "invalidhash"

	// Should return false for invalid hash
	if utils.CheckPassword(password, invalidHash) {
		t.Error("Expected password to not match invalid hash")
	}
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	password := ""

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if hashedPassword == "" {
		t.Error("Expected hashed password even for empty string")
	}
}
