package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveAndLoad(t *testing.T) {
	// Setup temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "koeda_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Mock XDG_CACHE_HOME
	t.Setenv("XDG_CACHE_HOME", tmpDir)

	repos := []Repo{
		{FullName: "user/repo1", Private: false, UpdatedAt: time.Now()},
		{FullName: "user/repo2", Private: true, UpdatedAt: time.Now()},
	}

	// Test Save
	if err := Save(repos); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	expectedPath := filepath.Join(tmpDir, "koeda", "repos.json")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Cache file was not created at %s", expectedPath)
	}

	// Test Load
	data, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(data.Repositories) != 2 {
		t.Errorf("Expected 2 repos, got %d", len(data.Repositories))
	}

	if data.Repositories[0].FullName != "user/repo1" {
		t.Errorf("Unexpected repo name: %s", data.Repositories[0].FullName)
	}
}

func TestIsExpired(t *testing.T) {
	now := time.Now()
	ttl := 1 * time.Hour

	tests := []struct {
		name       string
		lastUpdate time.Time
		want       bool
	}{
		{
			name:       "Not Expired",
			lastUpdate: now.Add(-30 * time.Minute),
			want:       false,
		},
		{
			name:       "Expired",
			lastUpdate: now.Add(-2 * time.Hour),
			want:       true,
		},
		{
			name:       "Just Expired", // Depending on implementation detail, usually > TTL
			lastUpdate: now.Add(-1 * time.Hour).Add(-1 * time.Second),
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsExpired(tt.lastUpdate, ttl); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCacheDir(t *testing.T) {
	// Test with XDG_CACHE_HOME
	tmpDir := "/tmp/xdg_cache"
	t.Setenv("XDG_CACHE_HOME", tmpDir)

	dir, err := GetCacheDir()
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join(tmpDir, "koeda")
	if dir != expected {
		t.Errorf("Expected %s, got %s", expected, dir)
	}

	// Test fallback to HOME (clearing XDG_CACHE_HOME)
	t.Setenv("XDG_CACHE_HOME", "")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	expectedHome := filepath.Join(homeDir, ".cache", "koeda")

	dir, err = GetCacheDir()
	if err != nil {
		t.Fatal(err)
	}
	if dir != expectedHome {
		t.Errorf("Expected %s, got %s", expectedHome, dir)
	}
}
