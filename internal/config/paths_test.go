package config

import (
	"os/user"
	"testing"
)

func TestGetCSheetDir(t *testing.T) {
	result := GetCSheetDir()

	if result == "" {
		t.Error("expected non-empty directory path")
	}

	// Verify it matches the current user's home directory
	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("failed to get current user: %v", err)
	}

	if result != currentUser.HomeDir {
		t.Errorf("expected %q, got %q", currentUser.HomeDir, result)
	}
}
