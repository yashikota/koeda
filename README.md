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
- `--ttl`: Specifies the cache time-to-live and saves it to the configuration (default: `24h`)

## Configuration

Settings and cache are stored in the following paths (XDG Base Directory compliant):

- Cache: `~/.cache/koeda/repos.json` (or `$XDG_CACHE_HOME/koeda/repos.json`)
- Config: `~/.config/koeda/config.json` (or `$XDG_CONFIG_HOME/koeda/config.json`)

You can change the cache expiration time (TTL) by running `koeda update --ttl <duration>` or by editing `config.json`.
