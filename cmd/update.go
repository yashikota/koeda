package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/yashikota/koeda/internal/cache"
	"github.com/yashikota/koeda/internal/config"
	"github.com/yashikota/koeda/internal/github"
)

var UpdateCommand = &cli.Command{
	Name:  "update",
	Usage: "Update repository cache",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "affiliation",
			Usage: "Comma-separated list of affiliations (owner, collaborator, organization_member)",
			Value: "owner,collaborator,organization_member",
		},
		&cli.StringFlag{
			Name:  "visibility",
			Usage: "Visibility of repositories (all, public, private)",
			Value: "all",
		},
		&cli.DurationFlag{
			Name:  "ttl",
			Usage: "Cache time-to-live",
		},
	},
	Action: func(c *cli.Context) error {
		if c.IsSet("ttl") {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			cfg.TTL = c.Duration("ttl")
			if err := config.Save(cfg); err != nil {
				return err
			}
		}

		start := time.Now()
		count, err := DoUpdate(c.Context, github.FetchOptions{
			Affiliation: c.String("affiliation"),
			Visibility:  c.String("visibility"),
		})
		if err != nil {
			return err
		}
		fmt.Printf("Updated %d repositories in %v.\n", count, time.Since(start).Round(time.Millisecond))
		return nil
	},
}

// DoUpdate fetches repos and saves them to cache. Returns the number of repos fetched.
func DoUpdate(ctx context.Context, opts github.FetchOptions) (int, error) {
	client, err := github.NewClient(ctx)
	if err != nil {
		return 0, err
	}

	repos, err := client.FetchRepos(ctx, opts)
	if err != nil {
		return 0, err
	}

	if err := cache.Save(repos); err != nil {
		return 0, fmt.Errorf("failed to save cache: %w", err)
	}

	return len(repos), nil
}
