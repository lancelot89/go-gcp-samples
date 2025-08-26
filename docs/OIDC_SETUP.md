# GitHub Actions OIDC Setup for GCP

このドキュメントでは、GitHub ActionsとGCP間でOIDC (OpenID Connect) を使用した認証を設定する手順を説明します。この設定により、サービスアカウントキーを使わずに安全にGCPリソースへアクセスできます。

## 目次

1. [概要](#概要)
2. [前提条件](#前提条件)
3. [GCP側の設定](#gcp側の設定)
4. [GitHub側の設定](#github側の設定)
5. [動作確認](#動作確認)
6. [トラブルシューティング](#トラブルシューティング)
7. [セキュリティベストプラクティス](#セキュリティベストプラクティス)

## 概要

Workload Identity Federation を使用することで：
- 🔐 サービスアカウントキーの管理が不要
- 🚀 短期間の認証トークンを使用してセキュリティを向上
- 🔄 キーローテーション不要
- 📊 CloudAudit Logsで詳細な監査が可能

## 前提条件

- GCPプロジェクトが作成済み
- 必要なAPIが有効化済み：
  - IAM Service Account Credentials API
  - Security Token Service API
- `gcloud` CLIがインストール済み
- プロジェクトオーナーまたはIAM管理者権限

## GCP側の設定

### 1. 環境変数の設定

```bash
# プロジェクトIDを設定
export PROJECT_ID="your-gcp-project-id"

# その他の変数
export POOL_NAME="github-actions-pool"
export POOL_DISPLAY_NAME="GitHub Actions Pool"
export PROVIDER_NAME="github-provider"
export SA_NAME="github-actions-sa"
export SA_DISPLAY_NAME="GitHub Actions Service Account"
export GITHUB_REPO="lancelot89/go-gcp-samples"  # owner/repo形式

# gcloudのデフォルトプロジェクトを設定
gcloud config set project ${PROJECT_ID}
```

### 2. Workload Identity Poolの作成

```bash
gcloud iam workload-identity-pools create ${POOL_NAME} \
    --location="global" \
    --display-name="${POOL_DISPLAY_NAME}" \
    --description="Identity pool for GitHub Actions OIDC"
```

### 3. Workload Identity Providerの作成

```bash
gcloud iam workload-identity-pools providers create-oidc ${PROVIDER_NAME} \
    --location="global" \
    --workload-identity-pool=${POOL_NAME} \
    --display-name="GitHub Provider" \
    --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository,attribute.repository_owner=assertion.repository_owner" \
    --attribute-condition="assertion.repository == '${GITHUB_REPO}'" \
    --issuer-uri="https://token.actions.githubusercontent.com"
```

> **注意**: `--attribute-condition` により、指定したリポジトリからのアクセスのみが許可されます。

### 4. サービスアカウントの作成

```bash
# サービスアカウントを作成
gcloud iam service-accounts create ${SA_NAME} \
    --display-name="${SA_DISPLAY_NAME}" \
    --description="Service account for GitHub Actions CI/CD"

# 必要な権限を付与
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/artifactregistry.writer"

gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/run.developer"

gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountUser"
```

### 5. Workload Identity PoolとService Accountの紐付け

```bash
gcloud iam service-accounts add-iam-policy-binding \
    ${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com \
    --role="roles/iam.workloadIdentityUser" \
    --member="principalSet://iam.googleapis.com/projects/$(gcloud projects describe ${PROJECT_ID} --format='value(projectNumber)')/locations/global/workloadIdentityPools/${POOL_NAME}/attribute.repository/${GITHUB_REPO}"
```

### 6. 設定値の確認

```bash
# Workload Identity Provider のリソース名を取得
WIF_PROVIDER=$(gcloud iam workload-identity-pools providers describe ${PROVIDER_NAME} \
    --workload-identity-pool=${POOL_NAME} \
    --location=global \
    --format="value(name)")

# Service Account のメールアドレスを取得
WIF_SERVICE_ACCOUNT="${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

echo "========================================"
echo "GitHub Secrets に設定する値:"
echo "========================================"
echo "PROJECT_ID: ${PROJECT_ID}"
echo "REGION: asia-northeast1"
echo "WIF_PROVIDER: ${WIF_PROVIDER}"
echo "WIF_SERVICE_ACCOUNT: ${WIF_SERVICE_ACCOUNT}"
echo "========================================"
```

## GitHub側の設定

### 1. リポジトリのSecrets設定

GitHubリポジトリの Settings > Secrets and variables > Actions から以下を追加：

| Secret名 | 値 | 説明 |
|---------|-----|------|
| `PROJECT_ID` | 上記で確認した値 | GCPプロジェクトID |
| `REGION` | `asia-northeast1` | デプロイ先リージョン |
| `WIF_PROVIDER` | 上記で確認した値 | Workload Identity Provider名 |
| `WIF_SERVICE_ACCOUNT` | 上記で確認した値 | サービスアカウントメール |

### 2. Artifact Registry リポジトリの作成

```bash
# Artifact Registry リポジトリを作成（未作成の場合）
gcloud artifacts repositories create samples \
    --repository-format=docker \
    --location=asia-northeast1 \
    --description="Docker repository for sample applications"
```

## 動作確認

### 1. GitHub Actions ワークフローの実行

1. GitHubリポジトリのActionsタブを開く
2. "Build and Push to GAR" ワークフローを選択
3. "Run workflow" をクリック
4. モジュールを選択して実行

### 2. ログの確認

正常に認証されていることを確認：
- "Authenticate to Google Cloud" ステップが成功
- Docker imageのpushが成功

### 3. GCPコンソールでの確認

Artifact Registryでイメージが作成されていることを確認：
```bash
gcloud artifacts docker images list \
    asia-northeast1-docker.pkg.dev/${PROJECT_ID}/samples
```

## トラブルシューティング

### 認証エラーが発生する場合

1. **attribute conditionの確認**
   ```bash
   # Provider の設定を確認
   gcloud iam workload-identity-pools providers describe ${PROVIDER_NAME} \
       --workload-identity-pool=${POOL_NAME} \
       --location=global
   ```

2. **IAM bindingの確認**
   ```bash
   # Service Account の IAM Policy を確認
   gcloud iam service-accounts get-iam-policy \
       ${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com
   ```

3. **Cloud Logging でエラーログを確認**
   ```bash
   gcloud logging read "resource.type=iam_service_account" \
       --limit=10 \
       --format=json
   ```

### よくあるエラーと対処法

| エラー | 原因 | 対処法 |
|-------|------|--------|
| `Permission 'iam.serviceAccounts.getAccessToken' denied` | Workload Identity User権限がない | IAM bindingを再設定 |
| `Attribute mapping error` | attribute-mappingが不正 | Provider設定を修正 |
| `Repository not allowed` | attribute-conditionが不一致 | リポジトリ名を確認 |

## セキュリティベストプラクティス

### 1. 最小権限の原則

サービスアカウントには必要最小限の権限のみを付与：
```bash
# 権限を確認
gcloud projects get-iam-policy ${PROJECT_ID} \
    --flatten="bindings[].members" \
    --filter="bindings.members:${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
    --format="table(bindings.role)"
```

### 2. Attribute Conditionの活用

特定のブランチやタグからのみアクセスを許可：
```bash
# mainブランチからのみ許可する例
--attribute-condition="assertion.repository == '${GITHUB_REPO}' && assertion.ref == 'refs/heads/main'"
```

### 3. 監査ログの有効化

Cloud Audit Logsで認証アクティビティを監視：
```bash
gcloud logging read "protoPayload.serviceName=sts.googleapis.com" \
    --limit=10 \
    --format=json
```

### 4. 定期的な権限レビュー

```bash
# 未使用のサービスアカウントを特定
gcloud recommender insights list \
    --location=global \
    --insight-type=google.iam.serviceAccount.Insight \
    --filter="stateInfo.state=ACTIVE"
```

## リソースのクリーンアップ

不要になった場合の削除手順：

```bash
# 1. Service Account の削除
gcloud iam service-accounts delete ${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com

# 2. Workload Identity Provider の削除
gcloud iam workload-identity-pools providers delete ${PROVIDER_NAME} \
    --workload-identity-pool=${POOL_NAME} \
    --location=global

# 3. Workload Identity Pool の削除
gcloud iam workload-identity-pools delete ${POOL_NAME} \
    --location=global

# 4. Artifact Registry リポジトリの削除（必要な場合）
gcloud artifacts repositories delete samples \
    --location=asia-northeast1
```

## 参考資料

- [Google Cloud - Workload Identity Federation](https://cloud.google.com/iam/docs/workload-identity-federation)
- [GitHub Docs - About security hardening with OpenID Connect](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect)
- [google-github-actions/auth](https://github.com/google-github-actions/auth)