class Store {
  final int id;
  final String name;
  final String address;
  final double latitude;
  final double longitude;
  final double? distance;

  Store({
    required this.id,
    required this.name,
    required this.address,
    required this.latitude,
    required this.longitude,
    this.distance,
  });

  factory Store.fromJson(Map<String, dynamic> json) {
    return Store(
      id: json['id'],
      name: json['name'],
      address: json['address'],
      latitude: json['latitude'].toDouble(),
      longitude: json['longitude'].toDouble(),
      distance: json['distance']?.toDouble(),
    );
  }
}
