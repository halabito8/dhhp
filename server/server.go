package server

import (
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
)

func Start() {
	env, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serverURL := env["SERVER_URL"]

	listener, err := net.Listen("tcp", serverURL)
	if err != nil {
		fmt.Println("Error:", err)
	}

	defer listener.Close()

	fmt.Println("Listening on:", serverURL)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error listening:", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	log.Println("Connection made")
	defer conn.Close()
}
