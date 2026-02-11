package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/yashikota/koeda/internal/cache"
	"github.com/yashikota/koeda/internal/finder"
	"github.com/yashikota/koeda/internal/github"
)

var RootCommand = &cli.Command{
	Name:  "koeda",
	Usage: "Fuzzy find GitHub repositories",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "force-update",
			Usage: "Force update cache before finding",
		},
		&cli.DurationFlag{
			Name:  "ttl",
			Usage: "Cache time-to-live",
			Value: 24 * time.Hour,
		},
	},
	Action: func(c *cli.Context) error {
		forceUpdate := c.Bool("force-update")
		ttl := c.Duration("ttl")

		var repos []cache.Repo
		var err error

		// Try to load cache if not forced
		if !forceUpdate {
			var data *cache.Data
			data, err = cache.Load()
			if err == nil {
				// Check expiration
				if cache.IsExpired(data.LastUpdate, ttl) {
					// Expired
					err = fmt.Errorf("cache expired")
				} else {
					repos = data.Repositories
				}
			}
		}

		// Update if forced, not found, or expired
		if forceUpdate || err != nil {
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
				// If update fails and we have no cache, that's fatal.
				// If we had an expired cache, we might want to use it as fallback?
				// The spec says "Cache missing & API fail -> exit 1".
				// It doesn't explicitly say "Expired & API fail -> fallback".
				// "If JSON broken -> Auto re-fetch".
				// We'll error out for now to be safe.
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
