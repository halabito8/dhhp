package server

import (
	"dhhp/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

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
	defer conn.Close()

	log.Println("Connection made")

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	log.Printf("Received: %s", buffer[:n])
	symbol := strings.TrimSpace(string(buffer[:n]))

	res, err := http.Get("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=" + symbol + "&interval=5min&apikey=UAWADI46W0EEK2PT")
	if err != nil {
		log.Printf("error making http request: %s\n", err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}

	var stockData models.StockResponse
	if err := json.Unmarshal(body, &stockData); err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return
	}

	fmt.Printf("Symbol: %s\n", stockData.MetaData.Symbol)
	fmt.Printf("Last Refreshed: %s\n", stockData.MetaData.LastRefreshed)

	for timestamp, data := range stockData.TimeSeries {
		fmt.Printf("\nTimestamp: %s\n", timestamp)
		fmt.Printf("Open: %.4f\n", data.Open)
		fmt.Printf("Close: %.4f\n", data.Close)
		fmt.Printf("Volume: %d\n", data.Volume)
	}

	response := fmt.Sprintf("[TICKER:%s][PRICE:%f][TIME:%s]", symbol, stockData.TimeSeries["2024-12-26 19:55:00"].Close, "2024-12-26 19:55:00")

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error sending response to client:", err)
		return
	}

	fmt.Println("Response sent to client:", response)

}
