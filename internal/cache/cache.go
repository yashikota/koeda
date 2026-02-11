package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Repo struct {
	FullName  string    `json:"full_name"`
	Private   bool      `json:"private"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Data struct {
	Repositories []Repo    `json:"repositories"`
	LastUpdate   time.Time `json:"last_update"`
}

func GetCacheDir() (string, error) {
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		cacheDir = filepath.Join(home, ".cache")
	}
	return filepath.Join(cacheDir, "koeda"), nil
}

func Load() (*Data, error) {
	dir, err := GetCacheDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "repos.json")

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data Data
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func Save(repos []Repo) error {
	dir, err := GetCacheDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path := filepath.Join(dir, "repos.json")
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	data := Data{
		Repositories: repos,
		LastUpdate:   time.Now(),
	}

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func IsExpired(lastUpdate time.Time, ttl time.Duration) bool {
	return time.Since(lastUpdate) > ttl
}

// Convert helper if needed, but for now we keep it simple
