# Go x GCP Monorepo Samples

このリポジトリは、Go言語とGoogle Cloud Platform (GCP) を組み合わせた様々なサンプルコードを管理するモノレポです。

## 主な機能

- 🚀 **Cloud Run** と **Firestore** のサンプル実装
- 🔐 **GitHub Actions OIDC** によるセキュアなCI/CD
- 📦 **Go Workspace** を使用したモノレポ管理
- 🛡️ **Distroless** イメージによる軽量・セキュアなコンテナ
- ✅ 包括的なテストとLintのCI

## ディレクトリ構成

各サブディレクトリが独立したGCPサービスのサンプルプロジェクトに対応しています。

```
/
├── v1-cloud-run/  # Google Cloud Run サンプル
├── v2-firestore/  # Google Cloud Firestore サンプル
└── ...            # 今後、他のGCPサービスサンプルを追加予定
```

### `v1-cloud-run`

[go-cloudrun-example](https://github.com/lancelot89/go-cloudrun-example) の履歴を完全に引き継いだ、Google Cloud Runのサンプルアプリケーションです。

### `v2-firestore`

Google Cloud Firestore を利用したサンプルアプリケーションです。

## Goワークスペース

本リポジトリはGoのワークスペース機能 (`go.work`) を利用しており、ルートディレクトリから各モジュールのコマンドを横断的に実行できます。

```bash
# 例: 全モジュールのテストを実行
go test ./...
```

## CI/CD

このリポジトリでは、GitHub ActionsとGoogle CloudのOIDCを利用した安全なCI/CDパイプラインを提供しています。

### ワークフロー

- **Go CI** (`go-ci.yaml`): PR時のLint、テスト、ビルド、セキュリティスキャン
- **Build and Push to GAR** (`build-push-gar.yaml`): DockerイメージのビルドとArtifact Registryへのプッシュ
- **Deploy to Cloud Run** (`deploy-cloudrun.yaml`): Cloud Runへのデプロイ

### セットアップ

OIDC認証の設定手順は [docs/OIDC_SETUP.md](./docs/OIDC_SETUP.md) を参照してください。

## 開発

```bash
# Go 1.24.6をインストール
go version  # go version go1.24.6 linux/amd64

# 依存関係の取得
make init

# フォーマット
make fmt

# Lint
make lint

# テスト
make test

# ビルド
make build
```

## 背景

このモノレポは、既存の `go-cloudrun-example` リポジトリのコミット履歴を `git filter-repo` を用いて保持したまま、サブディレクトリに統合することで構築されました。
