package finder

import (
	"strings"
	"testing"
	"time"

	"github.com/yashikota/koeda/internal/cache"
)

func TestMakePreview(t *testing.T) {
	updatedAt := time.Date(2023, 10, 1, 12, 34, 56, 0, time.UTC)

	tests := []struct {
		name string
		repo cache.Repo
		want []string // Substrings that must be present
	}{
		{
			name: "Public Repo",
			repo: cache.Repo{
				FullName:  "user/public-repo",
				Private:   false,
				UpdatedAt: updatedAt,
			},
			want: []string{
				"Name: user/public-repo",
				"Visibility: Public",
				"Updated: 2023-10-01 12:34:56",
			},
		},
		{
			name: "Private Repo",
			repo: cache.Repo{
				FullName:  "user/private-repo",
				Private:   true,
				UpdatedAt: updatedAt,
			},
			want: []string{
				"Name: user/private-repo",
				"Visibility: Private",
				"Updated: 2023-10-01 12:34:56",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakePreview(tt.repo)
			for _, w := range tt.want {
				if !strings.Contains(got, w) {
					t.Errorf("MakePreview() = %q, missing %q", got, w)
				}
			}
		})
	}
}
