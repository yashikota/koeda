# 🌿 Koeda

`koeda` は GitHub のリポジトリ一覧をローカルにキャッシュし、高速に検索・選択するための CLI ツールです。  

## インストール

[Relases](https://github.com/yashikota/koeda/releases/latest) から最新のバイナリをダウンロードするか、Go Installしてください  

```bash
go install github.com/yashikota/koeda@latest
```

## 使い方

### 1. 認証 (推奨)

GitHub API のレート制限を回避し、プライベートリポジトリにアクセスするために認証設定を推奨します。  
`koeda` は以下の順序でトークンを探します  

1. 環境変数 `GITHUB_TOKEN`
2. `gh` CLI の認証情報 (`gh auth token`)

```bash
gh auth login
```

※ 認証がない場合、パブリックリポジトリのみ取得可能で、APIレート制限が厳しくなります。  

### 2. 実行

```bash
koeda
```

- 初回実行時は自動的にリポジトリ一覧を取得・キャッシュします。
- 2回目以降はキャッシュを使用するため高速に起動します。
- 選択したリポジトリ名（例: `owner/repo`）が標準出力に出力されます。

パイプで他のコマンドと組み合わせると便利です  

```bash
# 選択したリポジトリをcloneする
gh repo clone $(koeda)

# 選択したリポジトリをブラウザで開く
gh browse $(koeda)
```

### 3. キャッシュの更新

キャッシュを手動で更新する場合

```bash
koeda update
```

オプション
- `--visibility`: `all` (デフォルト), `public`, `private`
- `--affiliation`: `owner,collaborator,organization_member` (デフォルト)
- `--ttl`: キャッシュの有効期限を指定し、設定ファイルに保存します (デフォルト: `24h`)

## 設定

設定とキャッシュは以下のパスに保存されます  

- キャッシュ: `~/.cache/koeda/repos.json` (または `$XDG_CACHE_HOME/koeda/repos.json`)
- 設定: `~/.config/koeda/config.json` (または `$XDG_CONFIG_HOME/koeda/config.json`)

キャッシュの有効期限 (TTL) は、`koeda update --ttl <duration>` を実行するか、`config.json` を直接編集することで変更できます。
