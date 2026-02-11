package finder

import (
	"fmt"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/yashikota/koeda/internal/cache"
)

// MakePreview generates the preview string for a repository.
func MakePreview(r cache.Repo) string {
	vis := "Public"
	if r.Private {
		vis = "Private"
	}
	return fmt.Sprintf("Name: %s\nVisibility: %s\nUpdated: %s",
		r.FullName,
		vis,
		r.UpdatedAt.Format("2006-01-02 15:04:05"),
	)
}

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
			return MakePreview(repos[i])
		}),
	)
	if err != nil {
		return nil, err
	}
	return &repos[idx], nil
}
