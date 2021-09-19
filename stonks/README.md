# Stonks App

A small emulation of a Stock exchange written in [Go](https://golang.org/) for the backend and [Flutter](https://flutter.dev/) for the app.

### Getting Started

This project is a created using the `flutter create` command.

_Note: You need to have the Android SDK and Flutter installed. Follow [this](https://flutter.dev/docs/get-started/install) guide to setup your environment for your specific OS._

Steps to build this app:
1. Clone the repo
2. In a terminal, run: `cd stonks`
3. Then run `flutter pub get` to get the dependency package
4. Configure your server's IP address in `lib/main.dart` file in the `serverAddress` constant near the top of the file.
    - Note: Do not use `localhost` as your address since your server is not running on the device the app runs on. (Either physical or emulated).
5. Either connect your phone with an USB cable or start up an emulator
    - If you connect the phone with an USB cable, ensure USB debugging is enabled.
6. Run `flutter run` to build and run the app on the device.
7. Enjoy