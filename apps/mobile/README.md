# Flutter Mobile App

## API 接続先

API のベース URL は `--dart-define` で切り替え可能です。

```bash
# Android エミュレータ
flutter run --dart-define=API_BASE_URL=http://10.0.2.2:8080

# 実機 (PC の IP を指定)
flutter run --dart-define=API_BASE_URL=http://192.168.1.100:8080
```

## テスト実行

### 基本的なテスト
```bash
# すべてのテストを実行
flutter test

# 詳細モード
flutter test -v

# 特定のファイルのみ
flutter test test/models/store_test.dart

# Watch モード（変更を監視）
flutter test --watch
```

### テストカバレッジ
```bash
# カバレッジレポート生成
flutter test --coverage

# HTML レポート表示（lcov が必要）
brew install lcov  # macOS
genhtml coverage/lcov.info -o coverage/html
open coverage/html/index.html
```

### テスト構造

```
test/
├── models/              # モデルのユニットテスト
│   ├── store_test.dart
│   └── product_test.dart
├── services/            # API サービスのテスト
│   └── api_service_test.dart
├── screens/             # 画面のウィジェットテスト
│   └── home_screen_test.dart
└── widget_test.dart     # メインアプリのテスト
```

## テストの種類

### 1. ユニットテスト
モデルやビジネスロジックのテスト。

```dart
test('Store.fromJson should parse JSON correctly', () {
  final json = {'id': 1, 'name': 'テスト店舗'};
  final store = Store.fromJson(json);
  expect(store.id, 1);
});
```

### 2. ウィジェットテスト
UI コンポーネントのテスト。

```dart
testWidgets('App should display title', (tester) async {
  await tester.pumpWidget(const MyApp());
  expect(find.text('価格比較アプリ'), findsOneWidget);
});
```

### 3. 統合テスト（今後追加予定）
アプリ全体の動作テスト。

```bash
flutter drive --target=test_driver/app.dart
```

## Makefile コマンド

```bash
make test-mobile         # Flutter テスト実行
make test-mobile-watch   # Watch モード
```
