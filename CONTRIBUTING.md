# Contributing to go-gcp-samples

このドキュメントは、go-gcp-samplesプロジェクトへの貢献方法について説明します。

## 開発環境のセットアップ

### 必要なツール

- Go 1.24.6以上
- Docker（Cloud Runのローカルテスト用）
- gcloud CLI（GCPリソースの操作用）
- make（ビルドタスクの実行用）
- golangci-lint（Lintチェック用）

### 初期設定

```bash
# リポジトリのクローン
git clone https://github.com/lancelot89/go-gcp-samples.git
cd go-gcp-samples

# 依存関係のインストールとワークスペース初期化
make init

# Linterのインストール（必要な場合）
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## 開発フロー

### 1. Issueの確認

作業を始める前に、関連するIssueを確認するか、新しいIssueを作成してください。

### 2. ブランチの作成

```bash
# ブランチ命名規則: {type}/{issue-number}-{description}
git checkout -b feat/123-add-new-feature
git checkout -b fix/456-fix-bug
git checkout -b docs/789-update-readme
```

タイプの種類:
- `feat/`: 新機能
- `fix/`: バグ修正
- `docs/`: ドキュメント更新
- `chore/`: 雑務（依存関係更新等）
- `ci/`: CI/CD関連

### 3. 開発作業

#### コードスタイル

```bash
# フォーマット
make fmt

# Lintチェック
make lint

# フォーマット、Lint、テストをまとめて実行
make check
```

#### テスト

```bash
# 全テスト実行
make test

# 短縮テスト（エミュレータ不要）
make test-short
```

### 4. コミット

コミットメッセージは[Conventional Commits](https://www.conventionalcommits.org/)形式に従ってください。

```bash
# 例
git commit -m "feat(v1-cloud-run): add healthz endpoint"
git commit -m "fix(v2-firestore): fix connection timeout issue"
git commit -m "docs: update installation guide"
```

### 5. プルリクエスト

1. 変更をpush
```bash
git push origin feat/123-add-new-feature
```

2. GitHubでPRを作成
3. PRテンプレートに従って説明を記載
4. レビューを待つ

## コーディング規約

### Go

- 標準的なGoのコーディング規約に従う
- `go fmt`でフォーマット済みであること
- エラーは適切にハンドリングし、`errors.Is/As`を使用
- テストカバレッジ80%以上を目標

### ドキュメント

- 各モジュールにREADME.mdを必須
- 公開APIにはGoDocコメントを記載
- 複雑なロジックには説明コメントを追加

### セキュリティ

- APIキー、認証情報をコミットしない
- 環境変数で設定値を管理
- 最小権限の原則に従う

## モノレポ構造

```
go-gcp-samples/
├── v1-cloud-run/      # Cloud Runサンプル
├── v2-firestore/      # Firestoreサンプル
├── v3-bigquery/       # BigQueryサンプル（予定）
├── go.work           # Go Workspace設定
└── go.work.sum       # Workspace依存関係
```

新しいサンプルを追加する場合：
1. `vX-{service-name}/`形式でディレクトリ作成
2. `go work use ./vX-{service-name}`でワークスペースに追加
3. 独立したgo.modで管理

## 質問・サポート

- Issueで質問を投稿
- DiscussionsでアイデアやRFCを議論

## ライセンス

貢献いただいたコードは、本プロジェクトのライセンスに従います。