package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("Package Main:")
}

func main() {
	port := flag.Uint("port", 5000, "TCP port number for blockchain server")
	flag.Parse()
	app := NewBlockchainServer(uint16(*port))
	app.Run()
}
