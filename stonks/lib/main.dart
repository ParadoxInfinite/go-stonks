import 'dart:convert';

import 'package:flutter/foundation.dart';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:flutter/material.dart';

import 'models/stock.dart';

void main() => runApp(const MyApp());

const serverAddress = 'ws://192.168.88.58';

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    const title = 'Stonks';
    return const MaterialApp(
      title: title,
      home: MyHomePage(
        title: title,
      ),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({
    Key? key,
    required this.title,
  }) : super(key: key);

  final String title;

  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  final _channel = WebSocketChannel.connect(
    Uri.parse(serverAddress),
  );

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
      ),
      body: StreamBuilder(
        stream: _channel.stream.asBroadcastStream(),
        builder: (context, AsyncSnapshot snapshot) {
          if (snapshot.hasError) {
            return const SizedBox();
          }

          if (snapshot.connectionState == ConnectionState.waiting) {
            return const CircularProgressIndicator();
          }

          var json = jsonDecode(snapshot.data);
          List<Stock> stocks = [];
          // The JSON data is an array, so the decoded json is a list.
          // We will do the loop through this list to parse objects.
          if (json != null) {
            json.forEach((element) {
              stocks.add(Stock.fromJson(element));
            });
          }
          return ListView.separated(
            itemCount: stocks.length,
            separatorBuilder: (BuildContext context, int index) =>
                const Divider(height: 1),
            itemBuilder: (context, i) {
              final Stock stock = stocks[i];
              return StockTile(stock: stock);
            },
          );
        },
      ),
    );
  }

  @override
  void dispose() {
    _channel.sink.close();
    super.dispose();
  }
}

class StockTile extends StatelessWidget {
  final Stock stock;

  const StockTile({Key? key, required this.stock}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return ListTile(
      title: Text(stock.name),
      subtitle: Row(
        children: [
          Chip(
            label: Text.rich(TextSpan(
              text: stock.exchange,
              style: TextStyle(color: Colors.black.withOpacity(0.5)),
            )),
            padding: const EdgeInsets.fromLTRB(5.0, -3.0, 5.0, -3.0),
            labelPadding: const EdgeInsets.all(0),
          ),
          const Padding(padding: EdgeInsets.fromLTRB(16.0, 0, 0, 0)),
          Text.rich(TextSpan(
            text: stock.symbol,
            style: TextStyle(color: Colors.black.withOpacity(0.5)),
          ))
        ],
      ),
      trailing: Text(stock.price.toStringAsFixed(2)),
    );
  }
}
