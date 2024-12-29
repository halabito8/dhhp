package client

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

	conn, err := net.Dial("tcp", serverURL)
	if err != nil {
		log.Println("Error sending tcp:", err)
		return
	}

	defer conn.Close()

	data := []byte("MSFT")
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	log.Printf("Received: %s", buffer[:n])

}
