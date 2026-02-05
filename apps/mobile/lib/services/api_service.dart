import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/store.dart';
import '../models/product.dart';

class ApiService {
  // Change this to your computer's local IP address when testing on a real device
  // For emulator: use 10.0.2.2
  // For real device: use your computer's IP (e.g., 192.168.1.100)
  static const String baseUrl = 'http://localhost:8080';

  Future<List<Store>> getAllStores() async {
    try {
      final response = await http.get(Uri.parse('$baseUrl/api/stores'));

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        final List<dynamic> storesJson = data['stores'] ?? [];
        return storesJson.map((json) => Store.fromJson(json)).toList();
      } else {
        throw Exception('Failed to load stores');
      }
    } catch (e) {
      throw Exception('Failed to connect to server: $e');
    }
  }

  Future<List<Store>> getNearbyStores(double lat, double lon, {int radius = 5000}) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/api/stores/nearby?lat=$lat&lon=$lon&radius=$radius'),
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        final List<dynamic> storesJson = data['stores'] ?? [];
        return storesJson.map((json) => Store.fromJson(json)).toList();
      } else {
        throw Exception('Failed to load nearby stores');
      }
    } catch (e) {
      throw Exception('Failed to connect to server: $e');
    }
  }

  Future<List<Product>> getAllProducts() async {
    try {
      final response = await http.get(Uri.parse('$baseUrl/api/products'));

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        final List<dynamic> productsJson = data['products'] ?? [];
        return productsJson.map((json) => Product.fromJson(json)).toList();
      } else {
        throw Exception('Failed to load products');
      }
    } catch (e) {
      throw Exception('Failed to connect to server: $e');
    }
  }
}
