# 🌿 Koeda

`koeda` is a CLI tool that caches a list of GitHub repositories locally for fast searching and selection.

## Installation

Download the latest binary from [Releases](https://github.com/yashikota/koeda/releases/latest) or use Go Install:

```bash
go install github.com/yashikota/koeda@latest
```

## Usage

### 1. Authentication (Recommended)

Authentication is recommended to avoid GitHub API rate limits and to access private repositories.
`koeda` looks for tokens in the following order:

1. Environment variable `GITHUB_TOKEN`
2. `gh` CLI credentials (`gh auth token`)

```bash
gh auth login
```

*Note: Without authentication, only public repositories can be fetched, and stricter API rate limits will apply.*

### 2. Execution

```bash
koeda
```

- On the first run, it automatically fetches and caches the repository list.
- Subsequent runs use the cache for instant startup.
- The selected repository name (e.g., `owner/repo`) is printed to standard output.

It is useful to combine with other commands using pipes:

```bash
# Clone the selected repository
gh repo clone $(koeda)

# Open the selected repository in the browser
gh browse $(koeda)
```

### 3. Updating Cache

To manually update the cache:

```bash
koeda update
```

Options:
- `--visibility`: `all` (default), `public`, `private`
- `--affiliation`: `owner,collaborator,organization_member` (default)

### 4. Other Options

Root command options:

* `--force-update`: Forces a cache update before searching, even if a cache exists.
* `--ttl`: Specifies the cache time-to-live (default: `24h`).

```bash
# Update if the cache is older than 1 hour
koeda --ttl 1h
```

## Configuration

The cache file is stored in the following path (XDG Base Directory compliant):

- `~/.cache/koeda/repos.json`
- Or `$XDG_CACHE_HOME/koeda/repos.json`
