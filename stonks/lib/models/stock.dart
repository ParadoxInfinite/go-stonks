class Stock {
  final String id;
  final String name;
  final String symbol;
  final double price;
  final String exchange;

  Stock(
      {required this.id,
      required this.name,
      required this.symbol,
      required this.price,
      required this.exchange});

  factory Stock.fromJson(Map<String, dynamic> json) {
    return Stock(
      id: json['id'],
      name: json['name'],
      symbol: json['symbol'],
      price: json['price'],
      exchange: json['exchange'],
    );
  }
}
