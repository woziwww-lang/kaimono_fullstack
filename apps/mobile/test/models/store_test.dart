import 'package:flutter_test/flutter_test.dart';
import 'package:price_comparison_mobile/models/store.dart';

void main() {
  group('Store Model Tests', () {
    test('Store.fromJson should parse JSON correctly', () {
      final json = {
        'id': 1,
        'name': 'セブンイレブン 渋谷店',
        'address': '東京都渋谷区道玄坂1-2-3',
        'latitude': 35.6595,
        'longitude': 139.7007,
        'distance': 1234.56,
      };

      final store = Store.fromJson(json);

      expect(store.id, 1);
      expect(store.name, 'セブンイレブン 渋谷店');
      expect(store.address, '東京都渋谷区道玄坂1-2-3');
      expect(store.latitude, 35.6595);
      expect(store.longitude, 139.7007);
      expect(store.distance, 1234.56);
    });

    test('Store.fromJson should handle null distance', () {
      final json = {
        'id': 2,
        'name': 'ファミリーマート',
        'address': '東京都新宿区',
        'latitude': 35.6938,
        'longitude': 139.7034,
      };

      final store = Store.fromJson(json);

      expect(store.id, 2);
      expect(store.distance, isNull);
    });
  });
}
