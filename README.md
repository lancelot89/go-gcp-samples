# Go x GCP Monorepo Samples

このリポジトリは、Go言語とGoogle Cloud Platform (GCP) を組み合わせた様々なサンプルコードを管理するモノレポです。

## ディレクトリ構成

各サブディレクトリが独立したGCPサービスのサンプルプロジェクトに対応しています。

```
/
├── v1-cloud-run/  # Google Cloud Run サンプル
└── ...            # 今後、他のGCPサービスサンプルを追加予定
```

### `v1-cloud-run`

[go-cloudrun-example](https://github.com/lancelot89/go-cloudrun-example) の履歴を完全に引き継いだ、Google Cloud Runのサンプルアプリケーションです。

## Goワークスペース

本リポジトリはGoのワークスペース機能 (`go.work`) を利用しており、ルートディレクトリから各モジュールのコマンドを横断的に実行できます。

```bash
# 例: 全モジュールのテストを実行
go test ./...
```

## 背景

このモノレポは、既存の `go-cloudrun-example` リポジトリのコミット履歴を `git filter-repo` を用いて保持したまま、サブディレクトリに統合することで構築されました。
