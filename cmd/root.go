package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/yashikota/koeda/internal/cache"
	"github.com/yashikota/koeda/internal/config"
	"github.com/yashikota/koeda/internal/finder"
	"github.com/yashikota/koeda/internal/github"
)

var RootCommand = &cli.Command{
	Name:  "koeda",
	Usage: "Fuzzy find GitHub repositories",
	Action: func(c *cli.Context) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		ttl := cfg.TTL

		var repos []cache.Repo

		// Try to load cache
		data, err := cache.Load()
		if err == nil {
			// Check expiration
			if cache.IsExpired(data.LastUpdate, ttl) {
				// Expired
				err = fmt.Errorf("cache expired")
			} else {
				repos = data.Repositories
			}
		}

		// Update if not found or expired
		if err != nil {
			// If it was an error other than "not found", we might want to log it, but for now we just update.
			// Default options for auto-update
			opts := github.FetchOptions{
				Affiliation: "owner,collaborator,organization_member",
				Visibility:  "all",
			}

			// Show a spinner or message since this takes time
			fmt.Fprintln(os.Stderr, "Updating repository cache...")

			_, err := DoUpdate(c.Context, opts)
			if err != nil {
				return fmt.Errorf("failed to update cache: %w", err)
			}

			// Reload to get the fresh data
			data, loadErr := cache.Load()
			if loadErr != nil {
				return fmt.Errorf("failed to load cache after update: %w", loadErr)
			}
			repos = data.Repositories
		}

		if len(repos) == 0 {
			return fmt.Errorf("no repositories found")
		}

		selected, err := finder.Find(repos)
		if err != nil {
			if err.Error() == "abort" {
				// fuzzyfinder returns "abort" on user cancel (Ctrl+C/Esc)
				// Spec says exit 130
				return cli.Exit("cancelled", 130)
			}
			return err
		}

		fmt.Println(selected.FullName)
		return nil
	},
}
