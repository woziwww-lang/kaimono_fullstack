import 'package:flutter/material.dart';
import '../models/store.dart';
import '../models/product.dart';
import '../services/api_service.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  final ApiService _apiService = ApiService();
  List<Store> _stores = [];
  List<Product> _products = [];
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final stores = await _apiService.getAllStores();
      final products = await _apiService.getAllProducts();

      setState(() {
        _stores = stores;
        _products = products;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  Future<void> _loadNearbyStores() async {
    // Using Tokyo Station coordinates as default
    const lat = 35.6812;
    const lon = 139.7671;

    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final stores = await _apiService.getNearbyStores(lat, lon);

      setState(() {
        _stores = stores;
        _isLoading = false;
      });

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('東京駅周辺の店舗を読み込みました')),
        );
      }
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('価格比較アプリ'),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: _loadData,
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _error != null
              ? Center(
                  child: Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        const Icon(Icons.error, size: 48, color: Colors.red),
                        const SizedBox(height: 16),
                        Text(
                          'エラー: $_error',
                          textAlign: TextAlign.center,
                          style: const TextStyle(color: Colors.red),
                        ),
                        const SizedBox(height: 16),
                        ElevatedButton(
                          onPressed: _loadData,
                          child: const Text('再読み込み'),
                        ),
                      ],
                    ),
                  ),
                )
              : ListView(
                  padding: const EdgeInsets.all(16),
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text(
                          '近くの店舗',
                          style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
                        ),
                        ElevatedButton(
                          onPressed: _loadNearbyStores,
                          child: const Text('東京駅周辺'),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),
                    ..._stores.map((store) => Card(
                          margin: const EdgeInsets.only(bottom: 12),
                          child: ListTile(
                            leading: const Icon(Icons.store, color: Colors.blue),
                            title: Text(store.name),
                            subtitle: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(store.address),
                                if (store.distance != null)
                                  Text(
                                    '距離: ${(store.distance! / 1000).toStringAsFixed(2)} km',
                                    style: const TextStyle(
                                      color: Colors.blue,
                                      fontWeight: FontWeight.bold,
                                    ),
                                  ),
                              ],
                            ),
                          ),
                        )),
                    const SizedBox(height: 24),
                    const Text(
                      '商品一覧',
                      style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
                    ),
                    const SizedBox(height: 16),
                    ..._products.map((product) => Card(
                          margin: const EdgeInsets.only(bottom: 12),
                          child: ListTile(
                            leading: const Icon(Icons.shopping_basket, color: Colors.green),
                            title: Text(product.name),
                            subtitle: Text('${product.category} • ${product.barcode}'),
                          ),
                        )),
                  ],
                ),
    );
  }
}
