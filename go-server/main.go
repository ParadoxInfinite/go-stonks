package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":80", "http service address")

func main() {
	flag.Parse()
	exch := newExchange()
	go exch.run()
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		serveAuth(w, r)
	})
	http.HandleFunc("/prices", func(w http.ResponseWriter, r *http.Request) {
		serveWs(exch, w, r)
	})
	log.Printf("Listening to websocket requests at %s", *addr)
	e := http.ListenAndServe(*addr, nil)
	if e != nil {
		log.Fatal("Error in ListenAndServe: ", e)
	}
}
