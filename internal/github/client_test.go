package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/v69/github"
)

func TestFetchRepos_Pagination(t *testing.T) {
	// Setup mock server
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	// Mock handler for /user/repos
	mux.HandleFunc("/user/repos", func(w http.ResponseWriter, r *http.Request) {
		// Check parameters
		if r.URL.Query().Get("per_page") != "100" {
			t.Errorf("Expected per_page=100")
		}
		if r.URL.Query().Get("visibility") != "all" {
			t.Errorf("Expected visibility=all")
		}

		page := r.URL.Query().Get("page")
		if page == "" || page == "1" {
			// Page 1
			w.Header().Set("Link", fmt.Sprintf(`<%s/user/repos?page=2>; rel="next", <%s/user/repos?page=2>; rel="last"`, server.URL, server.URL))
			fmt.Fprint(w, `[
				{"full_name": "owner/repo1", "private": false, "updated_at": "2023-01-01T00:00:00Z"},
				{"full_name": "owner/repo2", "private": true, "updated_at": "2023-01-02T00:00:00Z"}
			]`)
		} else if page == "2" {
			// Page 2 (Last)
			fmt.Fprint(w, `[
				{"full_name": "owner/repo3", "private": false, "updated_at": "2023-01-03T00:00:00Z"}
			]`)
		} else {
			t.Errorf("Unexpected page requested: %s", page)
		}
	})

	// Create client pointing to mock server
	// Since client.client is private, we construct it manually for testing.
	// This is possible because we are in package github.
	client := &Client{
		client: github.NewClient(nil),
	}
	url, _ := url.Parse(server.URL + "/")
	client.client.BaseURL = url
	client.client.UploadURL = url

	opts := FetchOptions{
		Visibility: "all",
	}

	repos, err := client.FetchRepos(context.Background(), opts)
	if err != nil {
		t.Fatalf("FetchRepos failed: %v", err)
	}

	if len(repos) != 3 {
		t.Errorf("Expected 3 repos, got %d", len(repos))
	}

	expectedNames := []string{"owner/repo1", "owner/repo2", "owner/repo3"}
	for i, name := range expectedNames {
		if repos[i].FullName != name {
			t.Errorf("Repo[%d] name mismatch: got %s, want %s", i, repos[i].FullName, name)
		}
	}

	// Check private flag mapping
	if repos[0].Private != false {
		t.Errorf("Repo 1 should be public")
	}
	if repos[1].Private != true {
		t.Errorf("Repo 2 should be private")
	}
}
