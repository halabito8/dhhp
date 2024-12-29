package main

import (
	"dhhp/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Be more restrictive in production
	},
}

func main() {
	http.HandleFunc("/", handleWebSocket)
	port := ":8080" // You can get this from env too
	fmt.Printf("Starting WebSocket server on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("New WebSocket connection established")

	for {
		// Read message from browser
		_, symbol, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Get stock data
		response := processStockRequest(string(symbol))

		// Send the response back to the browser
		err = conn.WriteMessage(websocket.TextMessage, []byte(response))
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func processStockRequest(symbol string) string {
	env, err := godotenv.Read("../.env")
	if err != nil {
		return fmt.Sprintf("[ERROR:Unable to load environment configuration]")
	}

	res, err := http.Get("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=" + symbol + "&interval=1min&apikey=" + env["API_KEY"])
	if err != nil {
		return fmt.Sprintf("[ERROR:Failed to fetch stock data]")
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return fmt.Sprintf("[ERROR:Failed to read response data]")
	}

	var stockData models.StockResponse
	if err := json.Unmarshal(body, &stockData); err != nil {
		return fmt.Sprintf("[ERROR:Failed to parse stock data]")
	}

	var firstTimestamp string
	var firstData models.TimeSeriesData
	for timestamp, data := range stockData.TimeSeries {
		firstTimestamp = timestamp
		firstData = data
		break
	}

	return fmt.Sprintf("[TICKER:%s][PRICE:%f][TIME:%s]",
		symbol,
		firstData.Close,
		firstTimestamp)
}
