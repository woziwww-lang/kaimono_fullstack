# Kaimono Web App

価格比較Webアプリケーション

## 📁 ディレクトリ構造

```
app/
├── components/
│   ├── features/        # 機能コンポーネント
│   │   ├── MapView.tsx           # 地図表示コンポーネント
│   │   └── PriceTrendChart.tsx   # 価格推移チャート
│   └── ui/             # UIコンポーネント（未使用）
├── lib/                # ユーティリティ関数
│   ├── api.ts          # API通信関数
│   └── format.ts       # フォーマット関数
├── hooks/              # カスタムフック（未使用）
├── (routes)/           # ルート定義
│   ├── page.tsx        # トップページ（店舗検索・地図）
│   └── stores/
│       └── [id]/
│           └── page.tsx # 店舗詳細ページ
├── globals.css         # グローバルCSS
└── layout.tsx          # ルートレイアウト
```

## 🚀 主な機能

### トップページ (`/`)
- 🔍 商品名検索
- 🗺️ 地図表示（Leaflet + Supercluster）
- 📍 現在地取得
- 💰 価格順ソート（デフォルト）
- 📏 距離順ソート
- 🏪 店舗リスト表示

### 店舗詳細ページ (`/stores/[id]`)
- 📊 過去14日間の価格推移チャート
- 💹 最低・最高・平均価格表示
- 🔎 商品名でフィルタリング
- 📍 店舗住所・電話番号表示

## 🎨 UIデザイン

- Tailwind CSS使用
- グラデーション・アニメーション効果
- レスポンシブデザイン
- 絵文字アイコンで視覚的にわかりやすく

## 🛠️ 技術スタック

- Next.js 14 (App Router)
- React 18
- TypeScript
- Tailwind CSS
- Leaflet (地図)
- react-window (仮想スクロール)
