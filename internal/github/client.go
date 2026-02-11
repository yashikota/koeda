package github

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/v69/github"
	"github.com/yashikota/koeda/internal/cache"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		// Try getting token from gh cli
		if path, err := exec.LookPath("gh"); err == nil {
			cmd := exec.CommandContext(ctx, path, "auth", "token")
			out, err := cmd.Output()
			if err == nil {
				token = strings.TrimSpace(string(out))
			}
		}
	}

	var tc *http.Client
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = oauth2.NewClient(ctx, ts)
	} else {
		fmt.Fprintln(os.Stderr, "Warning: No GitHub token found. Using unauthenticated client (public repos only, low rate limit).")
	}

	return &Client{
		client: github.NewClient(tc),
	}, nil
}

type FetchOptions struct {
	Affiliation string
	Visibility  string
}

func (c *Client) FetchRepos(ctx context.Context, opts FetchOptions) ([]cache.Repo, error) {
	var allRepos []cache.Repo
	
	// Set default options if empty
	if opts.Affiliation == "" {
		opts.Affiliation = "owner,collaborator,organization_member"
	}
	if opts.Visibility == "" {
		opts.Visibility = "all"
	}

	opt := &github.RepositoryListByAuthenticatedUserOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Sort:        "updated",
		Direction:   "desc",
		Affiliation: opts.Affiliation,
		Visibility:  opts.Visibility,
	}

	for {
		repos, resp, err := c.client.Repositories.ListByAuthenticatedUser(ctx, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories: %w", err)
		}

		for _, r := range repos {
			if r.FullName == nil || r.UpdatedAt == nil {
				continue
			}
			allRepos = append(allRepos, cache.Repo{
				FullName:  *r.FullName,
				Private:   r.GetPrivate(),
				UpdatedAt: r.UpdatedAt.Time,
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil
}
