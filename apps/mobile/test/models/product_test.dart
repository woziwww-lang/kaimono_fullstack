import 'package:flutter_test/flutter_test.dart';
import 'package:price_comparison_mobile/models/product.dart';

void main() {
  group('Product Model Tests', () {
    test('Product.fromJson should parse JSON correctly', () {
      final json = {
        'id': 1,
        'name': 'コカ・コーラ 500ml',
        'category': '飲料',
        'barcode': '4902102072706',
      };

      final product = Product.fromJson(json);

      expect(product.id, 1);
      expect(product.name, 'コカ・コーラ 500ml');
      expect(product.category, '飲料');
      expect(product.barcode, '4902102072706');
    });

    test('Product should have all required fields', () {
      final product = Product(
        id: 2,
        name: '白米 5kg',
        category: '食品',
        barcode: '4901010012345',
      );

      expect(product.id, isNotNull);
      expect(product.name, isNotEmpty);
      expect(product.category, isNotEmpty);
      expect(product.barcode, isNotEmpty);
    });
  });
}
