以下では **「go-gcp-samples」という** *まっさらな GitHub レポジトリ* を作り、
既存の **`go-cloudrun-example` の履歴をまるごと保ったまま** サブディレクトリ `v1-cloud-run/` に取り込む、という“モノレポ化（履歴統合）”の手順を具体的に示します。
── **`git filter-repo` を使う方法** が最もクリーンかつ履歴をフルで残せるので、そのやり方をメインに解説し、最後に `git subtree`（履歴を 1 コミットに潰す簡易版）も添えています。

---

## 0. 前提

| ツール                                                          | バージョンめやす | 備考                                 |
| ------------------------------------------------------------ | -------- | ---------------------------------- |
| Git                                                          | ≥ 2.25   | `filter-repo` は新しめの Git 推奨         |
| [git-filter-repo](https://github.com/newren/git-filter-repo) | 最新       | `pipx install git-filter-repo` が手軽 |
| GitHub CLI（任意）                                               | ≥ 2.0    | `gh repo create` で空レポを作ると楽         |

---

## 1. 空の新レポジトリを GitHub に用意

```bash
# 例：GitHub CLI を使う場合
gh repo create lancelot89/go-gcp-samples --public --description "Go×GCP mono-repo samples" --homepage "https://lancelot89.dev"
```

* **まだ何も push しない**（＝空レポのまま）。

---

## 2. 既存レポをローカルに clone

```bash
git clone https://github.com/lancelot89/go-cloudrun-example.git
cd go-cloudrun-example
```

---

## 3. 既存履歴をサブディレクトリ付きに書き換え

### (git filter-repo ― フル履歴保持版)

> **目的**：すべてのファイルを `v1-cloud-run/…` の下に移動したコミット列を *作り直す*
> （元ブランチは触らないので、作業ブランチで実施）

```bash
# 新しい作業ブランチを切る
git switch -c migrate/monorepo

# フィルター適用：--path-rename '':'v1-cloud-run/'
git filter-repo --to-subdirectory-filter v1-cloud-run
# ↑ 全コミットのファイルパスを v1-cloud-run/ 以下へリライト
```

* コミットハッシュは置き換わりますが、**日時・メッセージ・作者情報は保持**。
* タグがある場合は `--tag-rename '':'v1-'` などでぶつからないように変更可。

---

## 4. 新レポ (go-gcp-samples) を push 先に設定

```bash
# 新しい origin を追加
git remote add mono git@github.com:lancelot89/go-gcp-samples.git

# まず main ブランチを push
git push mono migrate/monorepo:main
```

> **結果**
> `go-gcp-samples` の `main` には
> `v1-cloud-run/` 以下に第1回のコードがあり、コミット履歴はそのまま残ります。

---

## 5. 新しいワークスペース構成ファイルを追加

```bash
# go.work をルートに追加
cat <<'EOF' > go.work
go 1.24

use (
    ./v1-cloud-run
    ./v2-firestore
)
EOF

git add go.work
git commit -m "chore: add go.work for mono-repo workspace"
git push mono HEAD:main
```

---

## 6. 第2回（Firestore 編）のディレクトリを作成

```bash
mkdir v2-firestore
cd v2-firestore
go mod init github.com/lancelot89/go-gcp-samples/v2-firestore
# …Canvas のコードを配置…
git add .
git commit -m "feat(v2): scaffold Firestore sample"
git push mono HEAD:main
```

---

## 7. CI/CD（GitHub Actions）をモノレポ対応に

`.github/workflows/ci.yml` 例（モジュールごとにテスト）:

```yaml
name: Go CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        module:
          - v1-cloud-run
          - v2-firestore
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - run: go test ./... 
        working-directory: ${{ matrix.module }}
```

---

## 8. （任意）旧レポジトリ README でリダイレクト案内

* `go-cloudrun-example` のトップに

  > **「本レポジトリは `go-gcp-samples` に統合しました」**
  > と記載し、アーカイブ (readonly) 設定にすると混乱が減ります。

---

## 9. もし `git subtree` でサクッと統合したい場合（履歴を 1 コミットに“潰す”）

```bash
# 新しい空レポを clone して移動
git clone git@github.com:lancelot89/go-gcp-samples.git
cd go-gcp-samples

# 既存レポを remote として追加
git remote add v1 https://github.com/lancelot89/go-cloudrun-example.git

# サブツリーとして取り込み（--squash で履歴を 1 コミット化）
git subtree add --prefix v1-cloud-run v1 main --squash

git push origin main
```

* コミットは 1 つにまとまるので履歴がシンプル・容量も軽い。
* **細かな blame が不要**／「とにかく動くソースだけあれば OK」ならこちらでも可。

---

### まとめ

1. **履歴を丸ごと残すなら `git filter-repo`**

   * `--to-subdirectory-filter <dir>` でサブフォルダ化
   * そのまま新レポに push
2. **履歴を圧縮しても良いなら `git subtree --squash`**

   * 作業が数コマンドで完結

どちらの方法でも、新しい **`go-gcp-samples`** は

```
/
├── v1-cloud-run/
└── v2-firestore/
```

のモノレポ構成になります。
実際にコマンドを走らせる際にエラーや疑問が出たら、遠慮なく相談してください！
