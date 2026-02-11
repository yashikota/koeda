package finder

import (
	"fmt"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/yashikota/koeda/internal/cache"
)

func Find(repos []cache.Repo) (*cache.Repo, error) {
	idx, err := fuzzyfinder.Find(
		repos,
		func(i int) string {
			return repos[i].FullName
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			r := repos[i]
			vis := "Public"
			if r.Private {
				vis = "Private"
			}
			return fmt.Sprintf("Name: %s\nVisibility: %s\nUpdated: %s",
				r.FullName,
				vis,
				r.UpdatedAt.Format("2006-01-02 15:04:05"),
			)
		}),
	)
	if err != nil {
		return nil, err
	}
	return &repos[idx], nil
}
