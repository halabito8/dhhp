package main

import (
	"dhhp/client"
	"dhhp/server"
	"flag"
	"fmt"
)

func main() {
	mode := flag.String("mode", "server", "Mode to run: 'server' or 'client'")
	flag.Parse()

	switch *mode {
	case "server":
		fmt.Println("Starting server...")
		server.Start() // Call the server start function
	case "client":
		fmt.Println("Starting client...")
		client.Start() // Call the client start function
	default:
		fmt.Println("Invalid mode. Use 'server' or 'client'")
	}
}
