# Stonks Server

A small emulation of a Stock exchange written in [Go](https://golang.org/) for the backend and [Flutter](https://flutter.dev/) for the app.

### Getting Started

_Note: You need to have Go installed. Follow [this](https://golang.org/doc/install) guide to setup your environment for your specific OS._

Steps to build this app:
1. Clone the repo
2. In a terminal, run: `cd go-server`
3. Then run `go get .` to get the dependency packages
4. Run `go run .` to run the server.
    - The build by default runs the server on port 80. You can configure it by changing the `addr` variable's `":80"` string to whatever port. Note: Do not add `localhost` before the port because then the server is not exposed to your network, but only to your local machine.
    - Another note: You might get asked to allow permissions as this will expose a port of your local machine to the network.
5. Check the websocket with either a client like Postman _OR_ The Stonks App in this repo
7. Enjoy