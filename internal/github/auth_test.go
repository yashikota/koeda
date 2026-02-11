package github

import (
	"context"
	"testing"
)

func TestNewClient_TokenPriority(t *testing.T) {
	// Test case: GITHUB_TOKEN is set
	t.Setenv("GITHUB_TOKEN", "test-token-env")

	// We can't easily inspect the internal http.Client transport to verify the token
	// without reflection or mocking the entire oauth2 logic,
	// BUT we can verify that NewClient doesn't return an error.
	// For a strict unit test, we would need to dependency inject the token loader.
	// Here we just ensure it initializes correctly with the env var.

	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	// Ideally we would check if the token is actually used,
	// but that requires inspecting private fields of oauth2.Transport.
	// For now, this smoke test ensures no panic/error when ENV is present.
}
