import 'dart:convert';

import 'package:flutter/foundation.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:flutter/material.dart';

import 'models/stock.dart';

Future main() async {
  await dotenv.load(fileName: '.env');
  runApp(const Stonks());
}

class Stonks extends StatelessWidget {
  const Stonks({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    const title = 'Stonks';
    return const MaterialApp(
      title: title,
      home: StonksPage(
        title: title,
      ),
    );
  }
}

class StonksPage extends StatefulWidget {
  const StonksPage({
    Key? key,
    required this.title,
  }) : super(key: key);

  final String title;

  @override
  _StonksPageState createState() => _StonksPageState();
}

class _StonksPageState extends State<StonksPage> {
  final _channel = WebSocketChannel.connect(
    Uri.parse(
        '${dotenv.get('SECURE') == 'true' ? dotenv.get('SECURE_PROTOCOL') : dotenv.get('PROTOCOL')}://${dotenv.get('HOST')}'),
  );

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
      ),
      body: StreamBuilder(
        // To only listen to a stream, and not following the bad documentation,
        // @see https://api.flutter.dev/flutter/dart-async/Stream/asBroadcastStream.html
        stream: _channel.stream.asBroadcastStream(),
        builder: (context, AsyncSnapshot snapshot) {
          if (snapshot.hasError) {
            return const SizedBox();
          }

          if (snapshot.connectionState == ConnectionState.waiting) {
            // TODO: Either make a custom progress indicator, or center the default one.
            return const CircularProgressIndicator();
          }

          var json = jsonDecode(snapshot.data);
          List<Stock> stocks = [];
          // The JSON data is an array, so the decoded json is a list.
          // We will do the loop through this list to parse objects into the model.
          if (json != null) {
            json.forEach((element) {
              stocks.add(Stock.fromJson(element));
            });
          }
          // separated so we can add separatorBuilder to add a divider between items.
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
    // Closing whenever 'activity' dies.
    _channel.sink.close();
    super.dispose();
  }
}

// Takes a Stock modeled object and returns a ListTile for that object.
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
