package client

import (
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

	conn, err := net.Dial("tcp", serverURL)

	if err != nil {
		log.Println("Error sending tcp:", err)
		return
	}

	defer conn.Close()
}
