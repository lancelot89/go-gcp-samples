# Go Workspace 運用ガイド

このドキュメントでは、`go-gcp-samples`プロジェクトにおけるGo Workspaceの運用方法について説明します。

## 目次

- [Go Workspaceとは](#go-workspaceとは)
- [なぜGo Workspaceを使うのか](#なぜgo-workspaceを使うのか)
- [基本的な使い方](#基本的な使い方)
- [新しいモジュールの追加](#新しいモジュールの追加)
- [依存関係の管理](#依存関係の管理)
- [開発ワークフロー](#開発ワークフロー)
- [トラブルシューティング](#トラブルシューティング)
- [ベストプラクティス](#ベストプラクティス)

## Go Workspaceとは

Go 1.18から導入された機能で、複数のGoモジュールを一つのワークスペースで管理できる仕組みです。`go.work`ファイルを使用して、複数のモジュールを横断的に開発できます。

### 主な特徴

- 🔗 複数モジュールの同時開発
- 📦 モジュール間の依存関係をローカルで解決
- 🚀 モノレポでの開発効率向上
- 🔄 go.modファイルの独立性を維持

## なぜGo Workspaceを使うのか

### モノレポのメリット

1. **コードの再利用**: 共通コードを簡単に共有
2. **一元管理**: CI/CD、Lint、フォーマットを統一
3. **原子的な変更**: 複数モジュールの変更を1つのPRで
4. **依存関係の可視化**: モジュール間の関係が明確

### 従来の問題点と解決

| 問題点 | Go Workspaceでの解決 |
|--------|---------------------|
| `replace`ディレクティブの管理が複雑 | `go.work`で一元管理 |
| モジュール間の開発が非効率 | ローカルモジュールを直接参照 |
| CIでの`replace`削除忘れ | `go.work`はコミット対象外も可 |

## 基本的な使い方

### ワークスペースの初期化

```bash
# ルートディレクトリで実行
go work init

# 既存モジュールを追加
go work use ./v1-cloud-run
go work use ./v2-firestore
```

### go.workファイルの構造

```go
go 1.24.6

use (
    ./v1-cloud-run
    ./v2-firestore
)
```

### コマンド実行

```bash
# ルートから全モジュールのテスト実行
go test ./...

# 特定モジュールのビルド
go build ./v1-cloud-run/...

# 依存関係の同期
go work sync
```

## 新しいモジュールの追加

### 手順

1. **ディレクトリ作成**
```bash
mkdir v3-bigquery
cd v3-bigquery
```

2. **モジュール初期化**
```bash
go mod init github.com/lancelot89/go-gcp-samples/v3-bigquery
```

3. **基本構造の作成**
```bash
mkdir -p cmd/server internal pkg
touch README.md Dockerfile .env.example
```

4. **ワークスペースに追加**
```bash
cd ..
go work use ./v3-bigquery
```

5. **go.workの確認**
```bash
cat go.work
# use節にv3-bigqueryが追加されているか確認
```

### ディレクトリ命名規則

```
vX-{service-name}/
```

- `X`: バージョン番号（追加順）
- `service-name`: GCPサービス名（小文字、ハイフン区切り）

例：
- `v1-cloud-run`
- `v2-firestore`
- `v3-bigquery`
- `v4-pubsub`

## 依存関係の管理

### モジュール間の依存

同一ワークスペース内のモジュールを参照する場合：

```go
// v2-firestore/main.go
import (
    "github.com/lancelot89/go-gcp-samples/v1-cloud-run/pkg/util"
)
```

開発中は`go.work`により自動的にローカルモジュールが使用されます。

### 外部依存関係の追加

```bash
# 特定モジュールに依存関係を追加
cd v1-cloud-run
go get github.com/gorilla/mux@latest

# ワークスペース全体の依存関係を同期
cd ..
go work sync
```

### 依存関係の更新

```bash
# 全モジュールの依存関係を更新
for dir in v*/; do
    echo "Updating $dir"
    (cd "$dir" && go get -u ./... && go mod tidy)
done

# ワークスペースを同期
go work sync
```

## 開発ワークフロー

### 1. 日常的な開発

```bash
# プロジェクトルートで作業
cd go-gcp-samples

# フォーマット
make fmt

# Lint
make lint

# テスト
make test

# 特定モジュールのみテスト
go test ./v1-cloud-run/...
```

### 2. クロスモジュール開発

複数モジュールにまたがる変更を行う場合：

```bash
# 1. featureブランチを作成
git checkout -b feat/cross-module-update

# 2. 各モジュールで変更を実施
# v1-cloud-run/pkg/util/helper.go を編集
# v2-firestore/main.go で上記を利用

# 3. ルートから全テスト実行
go test ./...

# 4. 変更をコミット
git add .
git commit -m "feat: add shared utility across modules"
```

### 3. CI/CD対応

```yaml
# .github/workflows/go-ci.yaml の例
- name: Setup Go Workspace
  run: |
    go work sync
    
- name: Test all modules
  run: |
    go test -race -v ./...
```

## トラブルシューティング

### よくある問題と解決方法

#### 1. "go.work.sum is out of sync"

```bash
# 解決方法
go work sync
```

#### 2. "cannot find module providing package"

```bash
# モジュールがワークスペースに追加されているか確認
go work edit -json | jq '.Use'

# 追加されていない場合
go work use ./missing-module
```

#### 3. IDEが認識しない

VSCodeの場合：
```json
// .vscode/settings.json
{
    "go.experimentalWorkspaceModule": true
}
```

GoLand/IntelliJ IDEAの場合：
- Settings > Go > Go Modules
- "Enable Go modules integration" をチェック
- "Index entire GOPATH" のチェックを外す

#### 4. CIでgo.workが無視される

```bash
# CI環境では明示的に有効化
export GOWORK=go.work
go test ./...
```

## ベストプラクティス

### 1. go.workのバージョン管理

```gitignore
# ローカル開発用のgo.workは除外しない
# go.work

# ただしgo.work.sumは必須
# go.work.sum
```

### 2. モジュール間の依存を最小化

```go
// ❌ 悪い例: 直接内部実装に依存
import "github.com/lancelot89/go-gcp-samples/v1-cloud-run/internal/service"

// ✅ 良い例: 公開パッケージを使用
import "github.com/lancelot89/go-gcp-samples/v1-cloud-run/pkg/api"
```

### 3. 共通コードの配置

```
# 共通ユーティリティは専用モジュールに
common/
├── go.mod
├── pkg/
│   ├── logger/
│   ├── errors/
│   └── middleware/
```

### 4. モジュール単位のテスト

```makefile
# Makefile
.PHONY: test-modules
test-modules:
	@for dir in v*/; do \
		echo "Testing $$dir"; \
		(cd $$dir && go test -race -cover ./...); \
	done
```

### 5. 依存関係の定期更新

```bash
# 月次で実行
make update-deps

# または手動で
go work sync
go mod tidy -e
```

## 高度な使い方

### replaceディレクティブとの併用

開発中の外部モジュールを一時的に参照：

```go
// go.work
go 1.24.6

use (
    ./v1-cloud-run
    ./v2-firestore
)

replace github.com/some/module => ../local-module
```

### プライベートモジュールの参照

```bash
# Git認証設定
go env -w GOPRIVATE=github.com/yourorg/*

# SSH使用の強制
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

### ワークスペースの検証

```bash
# 構文チェック
go work edit -json | jq empty

# 使用モジュールの一覧
go work edit -json | jq -r '.Use[].DiskPath'

# 依存関係グラフの生成
go mod graph
```

## まとめ

Go Workspaceを使用することで：

- ✅ モノレポでの開発効率が向上
- ✅ モジュール間の依存関係が明確に
- ✅ CI/CDの統一が容易
- ✅ コードの再利用が促進

適切に運用することで、大規模プロジェクトでもスケーラブルな開発が可能になります。

## 参考リンク

- [Go Workspace Mode (公式ドキュメント)](https://go.dev/doc/tutorial/workspaces)
- [go.work file reference](https://go.dev/ref/mod#go-work-file)
- [Go Modules Reference](https://go.dev/ref/mod)