package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("Wallet server: ")
}
func main() {
	port := flag.Uint("Port", 8080, "TCP port number for Wallet Server")
	gateway := flag.String("gateway", "http://127.0.0.1:5001", "Blockchain Gateway")
	flag.Parse()
	app := NewWalletServer(uint16(*port), *gateway)
	app.Run()
}
