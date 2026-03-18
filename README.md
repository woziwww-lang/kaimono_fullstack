# 🛒 価格比較アプリ (Price Comparison App)
<img width="1505" height="748" alt="image" src="https://github.com/user-attachments/assets/4c5ef2e0-7e39-47aa-a885-ce90a4891ff2" />

[![Tech Stack](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go)](https://go.dev/)
[![Next.js](https://img.shields.io/badge/Next.js-14-000000?logo=next.js)](https://nextjs.org/)
[![Flutter](https://img.shields.io/badge/Flutter-3.0-02569B?logo=flutter)](https://flutter.dev/)
[![pnpm](https://img.shields.io/badge/pnpm-8-F69220?logo=pnpm)](https://pnpm.io/)

---

## 📑 目次

- [特徴](#-特徴)
- [技術スタック](#️-技術スタック)
- [環境構築](#-環境構築)
- [クイックスタート](#-クイックスタート)
- [プロジェクト構造](#-プロジェクト構造)
- [API エンドポイント](#-api-エンドポイント)
- [開発ガイド](#️-開発ガイド)
- [デプロイ](#-デプロイ)

---

## ✨ 特徴

- 🗺️ **地理空間検索**: PostGIS で近くの店舗を高速検索
- 💰 **リアルタイム価格比較**: 複数店舗の価格を一目で比較
- 📱 **マルチプラットフォーム**: Web (Next.js) + Mobile (Flutter)
- ⚡ **高速開発体験**: Turbopack + Vitest + pnpm (従来の 3.5 倍速)
- 🏗️ **モノレポ管理**: Turborepo で効率的なコード共有
- 🔒 **型安全**: TypeScript + Go の強い型システム

---

## 🏗️ 技術スタック

### コア技術

| カテゴリ | 技術 | 理由 |
|---------|------|------|
| **Monorepo** | Turborepo + pnpm workspace | 高速ビルド・キャッシング |
| **Package Manager** | pnpm 8.15 | npm より 2-3 倍高速 |
| **Backend** | Go 1.21 + Gin Framework | 高性能・低メモリ |
| **Database** | PostgreSQL 15 + PostGIS 3.3 | 地理空間クエリ |
| **Web Frontend** | Next.js 14 (App Router) | React Server Components |
| **Bundler** | Turbopack | Webpack より 10 倍高速 |
| **Styling** | Tailwind CSS 3.3 | Utility-first CSS |
| **Testing (Web)** | Vitest 1.1 | Jest より 5-10 倍高速 |
| **Testing (Mobile)** | Flutter Test | Dart 組込みテストフレームワーク |
| **Testing (Backend)** | Go Testing | 標準ライブラリ |
| **Mobile** | Flutter 3.0+ | iOS/Android 単一コードベース |
| **Infrastructure** | AWS + Terraform | IaC によるインフラ管理 |

### パフォーマンス比較

**開発サーバー起動時間:**
```
従来 (npm + Webpack + Jest):  ~60秒
最適化 (pnpm + Turbopack + Vitest): ~17秒  ⚡ 3.5倍高速
```

### テストフレームワーク比較

| プラットフォーム | フレームワーク | コマンド | 特徴 |
|---------------|--------------|---------|------|
| **Web** | Vitest | `make test-web` | Jest より 5-10 倍高速 |
| **Mobile** | Flutter Test | `make test-mobile` | Hot Reload 対応 |
| **Backend** | Go Testing | `make test-go` | 並列実行サポート |
| **All** | - | `make test` | 全テスト一括実行 |

---

## 🚀 環境構築

### 1. 前提条件

#### ✅ 必須（Web アプリ動作に必要）
- **Node.js** 18+ (現在: v22.1.0 ✓)
- **pnpm** 8+
- **Go** 1.21+
- **Docker Desktop**

#### 🔧 オプション
- **Flutter** 3.0+ (モバイルアプリ開発時のみ)
- **Terraform** (AWS デプロイ時のみ)

### 2. 自動インストール（macOS）

```bash
# 一括インストールスクリプト実行
./setup-mac.sh
```

### 3. 手動インストール（macOS）

```bash
# Homebrew がない場合
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 必須ツール
brew install go                  # Go 1.21+
brew install --cask docker       # Docker Desktop
npm install -g pnpm              # pnpm

# オプション
brew install --cask flutter      # Flutter (移動端のみ)
brew install terraform           # Terraform (AWS デプロイのみ)
```

⚠️ **Docker Desktop インストール後、必ずアプリケーションを起動してください**

### 4. 環境確認

```bash
make help                  # 利用可能なコマンド確認
node --version            # Node.js 確認
pnpm --version            # pnpm 確認
go version                # Go 確認
docker --version          # Docker 確認
```

---

## 🎯 クイックスタート

### ステップ 1: 依存関係のインストール

```bash
make install
```

これにより以下が実行されます：
- pnpm install (Web 依存関係)
- go mod download (Go 依存関係)
- flutter pub get (Flutter 依存関係)

### ステップ 2: データベース起動

```bash
make docker-up
make db-migrate-up
```

PostgreSQL + PostGIS が起動し、マイグレーションとサンプルデータが投入されます。

### ステップ 3: アプリケーション起動

**3つの別々のターミナルで実行:**

```bash
# ターミナル 1: Go バックエンド
make server
# → http://localhost:8080

# ターミナル 2: Next.js Web (Turbopack 使用)
make web
# → http://localhost:3000

# ターミナル 3: Flutter モバイル (オプション)
make mobile
```

### ステップ 4: 動作確認

**API テスト:**
```bash
# ヘルスチェック
curl http://localhost:8080/health

# 近くの店舗検索（東京駅周辺 5km）
curl "http://localhost:8080/api/stores/nearby?lat=35.6812&lon=139.7671&radius=5000" | jq

# 全店舗取得
curl http://localhost:8080/api/stores | jq

# 商品一覧
curl http://localhost:8080/api/products | jq
```

**Web アプリ:**
1. ブラウザで http://localhost:3000 を開く
2. 「東京駅周辺で検索」ボタンをクリック
3. 距離順にソートされた店舗リストを確認

### ステップ 5: テスト実行

```bash
# すべてのテストを実行（推奨）
make test

# 個別実行
make test-go       # Go バックエンドのみ
make test-web      # Next.js Web のみ
make test-mobile   # Flutter モバイルのみ

# 詳細テスト
cd apps/web && pnpm test:ui        # Vitest UI モード
cd apps/mobile && flutter test -v   # Flutter 詳細モード
```

---

## 📁 プロジェクト構造

```
/kaimono
├── apps/
│   ├── server/              # Go バックエンド API
│   │   ├── cmd/            # エントリーポイント
│   │   ├── internal/       # ビジネスロジック
│   │   │   ├── domain/    # ドメインモデル
│   │   │   ├── repository/# データアクセス層
│   │   │   ├── usecase/   # ユースケース
│   │   │   └── handler/   # HTTP ハンドラー
│   │   └── api/           # OpenAPI 定義
│   ├── web/               # Next.js Web アプリ
│   │   ├── app/          # App Router ページ
│   │   ├── components/   # React コンポーネント
│   │   └── __tests__/    # Vitest テスト
│   └── mobile/           # Flutter モバイルアプリ
│       ├── lib/
│       │   ├── models/   # データモデル
│       │   ├── services/ # API サービス
│       │   └── screens/  # 画面
│       └── pubspec.yaml
├── packages/
│   ├── database/         # SQL マイグレーション
│   │   └── init.sql     # 初期データ
│   └── shared-configs/  # 共有設定
├── infrastructure/
│   └── terraform/       # AWS リソース定義
├── docker-compose.yml   # ローカル開発環境
├── Makefile            # 開発コマンド
├── turbo.json          # Turborepo 設定
├── pnpm-workspace.yaml # pnpm workspace
└── README.md           # このファイル
```

---

## 🌐 API エンドポイント

> **認証 (任意)**: `API_KEY` を設定した場合、`X-API-Key` ヘッダー または `Authorization: Bearer <token>` が必要です。

### 店舗 (Stores)

| Method | Endpoint | 説明 | パラメータ |
|--------|----------|------|-----------|
| `GET` | `/api/stores` | 全店舗取得 | `q`, `category`, `bbox`, `user_lat`, `user_lon`, `limit`, `offset`, `sort`, `order` |
| `GET` | `/api/stores/nearby` | 近くの店舗検索 | `lat`, `lon`, `radius`, `limit`, `offset` |
| `GET` | `/api/stores/:id` | 店舗詳細 | - |
| `GET` | `/api/stores/:id/prices` | 店舗別価格一覧 | `category`, `limit`, `offset`, `sort`, `order` |

**例: 近くの店舗検索**
```bash
GET /api/stores/nearby?lat=35.6812&lon=139.7671&radius=5000
```

レスポンス:
```json
{
  "data": [
    {
      "id": 1,
      "name": "セブンイレブン 渋谷店",
      "address": "東京都渋谷区道玄坂1-2-3",
      "latitude": 35.6595,
      "longitude": 139.7007,
      "distance": 1234.56
    }
  ],
  "meta": {
    "count": 5,
    "limit": 20,
    "offset": 0
  }
}
```

### 商品 (Products)

| Method | Endpoint | 説明 | パラメータ |
|--------|----------|------|-----------|
| `GET` | `/api/products` | 全商品取得 | `limit`, `offset`, `sort`, `order` |
| `GET` | `/api/products/categories` | カテゴリ一覧 | - |
| `GET` | `/api/products/search` | 商品検索 | `q` (keyword), `limit`, `offset` |
| `GET` | `/api/products/:id` | 商品詳細 | - |
| `GET` | `/api/products/:id/prices` | 価格比較 | `limit`, `offset`, `sort`, `order` |

**例: 商品価格比較**
```bash
GET /api/products/1/prices
```

レスポンス:
```json
{
  "data": [
    {
      "id": 1,
      "price": 115.00,
      "currency": "JPY",
      "store": {
        "id": 2,
        "name": "ファミリーマート 新宿店"
      }
    }
  ],
  "meta": {
    "count": 1,
    "limit": 20,
    "offset": 0
  }
}
```

### ヘルスチェック

| Method | Endpoint | 説明 |
|--------|----------|------|
| `GET` | `/health` | サーバーステータス |
| `GET` | `/metrics` | Prometheus メトリクス |

API スキーマは `packages/shared-configs/openapi.yaml` にあります。

---

## 🛠️ 開発ガイド

### よく使うコマンド

```bash
make help            # すべてのコマンド表示
make install         # 依存関係インストール
make docker-up       # Docker 起動
make docker-down     # Docker 停止
make docker-logs     # Docker ログ表示
make db-status       # データベース状態確認
make db-migrate-up   # DB マイグレーション適用
make db-migrate-down # DB マイグレーション 1 つ戻す
make db-migrate-version # マイグレーション状態確認
make server          # Go サーバー起動
make web             # Next.js 起動 (Turbopack)
make mobile          # Flutter 起動
make test            # 全テスト実行 (Go + Web + Mobile)
make test-go         # Go テストのみ
make test-web        # Web テストのみ
make test-mobile     # Mobile テストのみ
make test-mobile-watch  # Mobile テスト (watch mode)
make clean           # ビルド成果物削除
make reset           # 完全リセット
```

### 環境変数

**Backend (.env)**
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=password
DB_NAME=price_comparison
DB_SSLMODE=disable
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
CACHE_TTL_SECONDS=60
API_KEY=
CORS_ORIGINS=http://localhost:3000,http://localhost:3001
METRICS_ROUTE=/metrics
LOG_LEVEL=info
PORT=8080
MIGRATIONS_PATH=../../packages/database/migrations
```

**Web (.env.local)**
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### 新機能の追加

#### 1. 新しい API エンドポイント (Go)

```bash
# 1. ドメインモデルを定義
apps/server/internal/domain/models.go

# 2. リポジトリを作成
apps/server/internal/repository/your_repository.go

# 3. ユースケース（ビジネスロジック）を作成
apps/server/internal/usecase/your_usecase.go

# 4. ハンドラーを作成
apps/server/internal/handler/your_handler.go

# 5. main.go にルートを追加
apps/server/cmd/main.go
```

#### 2. 新しいページ (Next.js)

```bash
# App Router を使用
apps/web/app/your-page/page.tsx
```

#### 3. テストの追加

**Web (Vitest)**
```bash
# テストファイル作成
apps/web/__tests__/your-feature.test.tsx

# 実行
make test-web
# または
cd apps/web && pnpm test
```

**Mobile (Flutter)**
```bash
# テストファイル作成
apps/mobile/test/your-feature_test.dart

# 実行
make test-mobile
# または
cd apps/mobile && flutter test
```

**Backend (Go)**
```bash
# テストファイル作成
apps/server/internal/handler/your_handler_test.go

# 実行
make test-go
# または
cd apps/server && go test ./...
```

### データベース操作

```bash
# データベースに接続
docker-compose exec db psql -U admin -d price_comparison

# テーブル確認
\dt

# 店舗データ確認
SELECT name, address FROM stores;

# 地理空間クエリ例
SELECT name, ST_Distance(
  location,
  ST_GeographyFromText('POINT(139.7671 35.6812)')
) as distance
FROM stores
ORDER BY distance
LIMIT 5;
```

### トラブルシューティング

**ポート競合**
```bash
# 使用中のポート確認
lsof -i :8080
lsof -i :3000

# プロセス終了
kill -9 <PID>
```

**Docker エラー**
```bash
# コンテナ再起動
make docker-down
make docker-up

# ログ確認
make docker-logs
```

**pnpm エラー**
```bash
# キャッシュクリア
pnpm store prune

# 再インストール
rm -rf node_modules
pnpm install
```

---

## ☁️ デプロイ

### AWS デプロイ（準備中）

```bash
cd infrastructure/terraform

# 初期化
terraform init

# プラン確認
terraform plan

# デプロイ実行
terraform apply
```

デプロイ詳細は `infrastructure/terraform/README.md` を参照（今後追加予定）。
