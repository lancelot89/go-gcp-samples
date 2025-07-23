# v2-firestore: Go + Cloud Firestore サンプルアプリケーション

このプロジェクトは、Go 言語と Google Cloud Platform (GCP) の Cloud Firestore を組み合わせたサンプルアプリケーションです。`article.md` で解説されているクラウドネイティブ開発のプラクティスに基づき、Firestore の主要な概念（コレクション、ドキュメント、インデックス）と Go クライアントを用いた CRUD 操作、トランザクション、そしてローカルエミュレータでの開発・テスト方法を実践的に示します。

## プロジェクトの概要

本アプリケーションは、シンプルな Todo 管理 API を提供します。Clean Architecture の原則に従い、以下のレイヤーで構成されています。

- `cmd/server`: アプリケーションのエントリポイント（HTTP サーバーの起動）
- `internal/firestore`: Firestore とのデータアクセス層（Repository）およびデータモデルの定義
- `internal/service`: ビジネスロジック層（Todo の作成、取得など）
- `internal/handler`: HTTP リクエストの処理層（API エンドポイントの定義）
- `pkg/config`: アプリケーションの設定管理

## 動作方法

### 前提条件

- Go (1.22+)
- Docker / Docker Compose
- gcloud CLI (Firestore Emulator の利用に必要)

### 1. ローカル開発環境のセットアップ

Firestore エミュレータと Go アプリケーションを Docker Compose でまとめて起動できます。

1.  **`docker-compose.yml` の `PROJECT_ID` を更新:**
    `docker-compose.yml` ファイルを開き、`PROJECT_ID: your-gcp-project-id` の部分を実際の GCP プロジェクト ID に置き換えてください。

2.  **Docker Compose で起動:**
    `v2-firestore` ディレクトリに移動し、以下のコマンドを実行します。
    ```bash
    docker compose up --build
    ```
    これにより、Firestore エミュレータと Go アプリケーションが起動します。Go アプリケーションはホストの `8081` ポートでリッスンします。

### 2. API エンドポイントの利用

サーバーが起動したら、`http://localhost:8081` に対して以下の API エンドポイントを利用できます。

#### Todo の作成 (POST)

- **URL:** `/todos`
- **Method:** `POST`
- **Content-Type:** `application/json`
- **Body:**
    ```json
    {
      "userId": "user123",
      "title": "Buy groceries"
    }
    ```
- **例 (curl):**
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"userId": "user123", "title": "Buy groceries"}' http://localhost:8081/todos
    ```

#### Todo の取得 (GET)

- **URL:** `/todos?id={todoId}`
- **Method:** `GET`
- **Query Parameter:** `id` (作成時に自動生成される Todo ID)
- **例 (curl):**
    ```bash
    curl "http://localhost:8081/todos?id=user123-Buy%20groceries"
    ```
    （`id` は `userId-title` の形式で生成されます。URL エンコードに注意してください。）

### 3. Cloud Build を使用したデプロイ

このプロジェクトは、Cloud Build を使用して GCP の Cloud Run にデプロイできるように設定されています。

1.  **Cloud Build API の有効化:**
    Google Cloud Console で Cloud Build API を有効にします。

2.  **Artifact Registry の設定:**
    Docker イメージを保存するための Artifact Registry リポジトリを作成します。
    ```bash
    gcloud artifacts repositories create cloud-run-source-deploy --repository-format=docker --location=asia-northeast1 --description="Docker repository for Cloud Run source deployments"
    ```

3.  **Cloud Build の実行:**
    `v2-firestore` ディレクトリに移動し、以下のコマンドを実行します。
    ```bash
    gcloud builds submit --config cloudbuild.yaml .
    ```
    これにより、Cloud Build がトリガーされ、アプリケーションのビルド、テスト、コンテナイメージのプッシュ、Cloud Run へのデプロイが自動的に行われます。

### 4. GitHub Actions を使用した Docker 環境の検証

`.github/workflows/v2-firestore-docker-check.yml` に定義された GitHub Actions ワークフローは、`v2-firestore` ディレクトリ内のファイルが変更された場合に自動的に実行されます。

このワークフローは、Docker Compose を使用してローカルでアプリケーションと Firestore エミュレータを起動し、API エンドポイントが正しく応答するかを検証します。これにより、変更が Docker 環境で問題なく動作することを確認できます。
