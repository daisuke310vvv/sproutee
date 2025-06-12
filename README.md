# Sproutee

Git Worktreeを活用した開発環境管理CLIツール

## 概要

Sprouteeは、Gitリポジトリ内でworktreeを効率的に管理し、必要なファイルを自動でコピーするCLIツールです。開発時に複数のブランチで作業する際に、環境設定ファイルやプロジェクト固有のファイルを簡単にworktreeに複製できます。

## 機能要件

### 1. Git Worktree作成機能
- 現在のリポジトリで`git worktree add`を実行
- worktreeの保存先：`.git/sproutee-worktrees/`ディレクトリ内
- ディレクトリ名：`[指定名]_[ランダム文字列]`の形式

### 2. ファイルコピー機能
- worktree作成後、指定されたファイルを自動コピー
- コピー対象ファイルは設定ファイルで管理

### 3. 設定ファイル管理
- コピー対象ファイルのリストを設定ファイルに記載
- 設定ファイル形式：JSON

### 4. コマンドライン インターフェース
- worktree名をコマンド引数で指定可能
- 直感的なコマンド構造

## 実装手順

### Phase 1: プロジェクト基盤構築
1. **Golang開発環境セットアップ**
   - asdfでGo 1.24.4をインストール
   - `go mod init github.com/[username]/sproutee`でプロジェクト初期化
   - プロジェクト構造作成（cmd/, internal/, pkg/）
   - go.modファイル設定

2. **基本CLIフレームワーク**
   - Cobraライブラリ導入（github.com/spf13/cobra）
   - コマンドライン引数パース機能
   - ヘルプメッセージ表示
   - エラーハンドリング基盤

### Phase 2: 設定ファイル機能
1. **設定ファイル構造設計**
   - 設定ファイル形式決定
   - デフォルト設定ファイル作成
   - 設定ファイル読み込み機能

2. **設定ファイル管理**
   - 設定ファイル存在チェック
   - 設定ファイル検証機能
   - 設定ファイル作成コマンド

### Phase 3: Git Worktree機能
1. **Git操作基盤**
   - Gitリポジトリ検出
   - Git worktreeコマンド実行
   - エラーハンドリング

2. **Worktree作成機能**
   - ランダム文字列生成
   - ディレクトリ名生成ロジック
   - worktree作成コマンド実行

### Phase 4: ファイルコピー機能
1. **ファイル操作**
   - ファイル存在チェック
   - ファイルコピー機能
   - ディレクトリ構造保持

2. **バッチコピー機能**
   - 設定ファイルからファイルリスト読み込み
   - 複数ファイル一括コピー
   - コピー結果レポート

### Phase 5: 統合とテスト
1. **機能統合**
   - 全機能の連携テスト
   - エラーケース対応
   - ユーザビリティ改善

2. **ドキュメント整備**
   - 使用方法ドキュメント
   - 設定ファイル例
   - トラブルシューティング

## コマンド仕様（予定）

```bash
# 基本的な使用方法
sproutee create <worktree-name> [branch-name]

# 設定ファイル管理
sproutee config init    # 設定ファイル作成
sproutee config list    # 設定内容表示
sproutee config edit    # 設定ファイル編集

# その他
sproutee list           # 作成済みworktree一覧
sproutee clean          # 不要なworktree削除
sproutee help           # ヘルプ表示
```

## 設定ファイル例

```json
// sproutee.json
{
  "copy_files": [
    ".env",
    ".env.local",
    "docker-compose.yml",
    "package-lock.json",
    "yarn.lock",
    "Makefile",
    ".vscode/settings.json"
  ]
}
```

## 開発TODO

### 🏗️ プロジェクト基盤（Golang）
- [x] asdfでGo 1.24.4の開発環境セットアップ
- [x] .tool-versionsファイル作成（`golang 1.24.4`）
- [x] `go mod init`でプロジェクト初期化
- [x] 標準的なGolangプロジェクト構造作成
  - [x] `cmd/sproutee/main.go`（エントリーポイント）
  - [x] `internal/`（内部パッケージ）
  - [x] `pkg/`（外部公開パッケージ）
- [x] go.mod依存関係設定
- [x] .gitignore作成（Go用）
- [x] Cobraライブラリ導入（CLI framework）

### ⚙️ 設定ファイル機能
- [x] JSON設定ファイル構造設計
- [x] encoding/json 標準ライブラリを使用
- [x] 設定ファイル読み込み機能実装
- [x] 設定ファイル作成コマンド実装（`sproutee config init`）
- [x] 設定ファイル検証機能実装
- [x] デフォルト設定ファイル（sproutee.json）作成
- [x] 設定ファイルパス検索機能（カレントディレクトリ→親ディレクトリ）

### 🌿 Worktree管理機能
- [x] Gitリポジトリ検出機能
- [x] ランダム文字列生成機能
- [x] ディレクトリ名生成ロジック
- [x] `git worktree add`コマンド実行機能
- [x] `.git/sproutee-worktrees/`ディレクトリ管理
- [x] Worktree作成コマンド実装

### 📁 ファイルコピー機能
- [x] ファイル存在チェック機能
- [x] 単一ファイルコピー機能
- [x] ディレクトリ構造保持コピー機能
- [x] 設定ファイルからファイルリスト読み込み
- [x] 複数ファイル一括コピー機能
- [x] コピー結果レポート機能

### 🖥️ CLI機能
- [x] コマンドライン引数パース
- [x] `create`コマンド実装（worktree作成＋ファイルコピー完全統合）
- [x] `config`サブコマンド実装（`init`, `list`）
- [x] `list`コマンド実装
- [x] `clean`コマンド実装
- [x] ヘルプメッセージ実装
- [ ] エラーメッセージ改善

### 🧪 テスト・品質
- [x] ユニットテスト作成
- [ ] 統合テスト作成
- [ ] エラーハンドリングテスト
- [ ] パフォーマンステスト
- [ ] マルチプラットフォーム対応確認

### 📖 ドキュメント
- [ ] 使用方法ドキュメント作成
- [ ] 設定ファイル仕様書作成
- [ ] トラブルシューティングガイド
- [ ] インストールガイド作成
- [ ] 開発者向けドキュメント

### 🚀 リリース準備（Homebrew配布）
- [ ] Goビルド設定（マルチプラットフォーム対応）
- [ ] Homebrew Formulaファイル作成
- [ ] GitHub Releases設定
- [ ] Homebrew tap リポジトリ作成
- [ ] GoReleaser設定（自動ビルド・リリース）
- [ ] CI/CDパイプライン構築（GitHub Actions）
- [ ] リリースノート作成

## 技術仕様

### 開発言語・フレームワーク
- **言語**: Go 1.24.4 (asdfで管理)
- **CLIフレームワーク**: Cobra (github.com/spf13/cobra)
- **設定ファイル**: JSON (encoding/json標準ライブラリ)
- **ビルドツール**: GoReleaser
- **配布**: Homebrew

### プロジェクト構造
```
sproutee/
├── cmd/sproutee/        # メインエントリーポイント
│   └── main.go
├── internal/            # 内部パッケージ
│   ├── config/         # 設定ファイル管理
│   ├── worktree/       # worktree操作
│   └── copy/           # ファイルコピー機能
├── pkg/                # 外部公開パッケージ
├── testdata/           # テストデータ
├── go.mod              # Go依存関係
├── go.sum
├── .goreleaser.yml     # リリース設定
├── .tool-versions      # asdf用バージョン管理
└── sproutee.rb         # Homebrew Formula
```

### 開発メモ
- worktreeディレクトリ：`.git/sproutee-worktrees/[名前]_[ランダム文字列]/`
- 設定ファイル名：`sproutee.json`
- ランダム文字列：8文字程度の英数字
- エラーハンドリング：Gitコマンドエラー、ファイルアクセスエラー等を適切に処理
- Homebrew Formula名：`sproutee`

## インストール

### Homebrew（リリース後）
```bash
brew tap [username]/sproutee
brew install sproutee
```

### 開発版
```bash
# Go 1.24.4をasdfでインストール
asdf plugin add golang https://github.com/asdf-community/asdf-golang.git
asdf install golang 1.24.4
asdf global golang 1.24.4

# プロジェクトクローンしてビルド
git clone https://github.com/[username]/sproutee.git
cd sproutee
go build -o sproutee cmd/sproutee/main.go
```

## 使用例

```bash
# worktreeを作成してファイルをコピー
sproutee create feature-123 develop

# 設定ファイルを初期化
sproutee config init

# 作成済みworktreeを確認
sproutee list
```

## ライセンス

MIT License
