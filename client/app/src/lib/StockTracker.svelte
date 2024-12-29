<script>
    import { onMount, onDestroy } from 'svelte';

    let socket;
    let symbol = '';
    let connected = false;
    let stockData = null;
    let error = null;

    function connectWebSocket() {
        // Replace with your WebSocket server URL
        socket = new WebSocket('ws://localhost:8080');

        socket.onopen = () => {
            connected = true;
            console.log('Connected to WebSocket');
        };

        socket.onmessage = (event) => {
            try {
                // Parse your custom protocol format [TICKER:XXX][PRICE:XXX][TIME:XXX]
                const data = event.data;
                const tickerMatch = data.match(/\[TICKER:(.*?)\]/);
                const priceMatch = data.match(/\[PRICE:(.*?)\]/);
                const timeMatch = data.match(/\[TIME:(.*?)\]/);

                if (tickerMatch && priceMatch && timeMatch) {
                    stockData = {
                        ticker: tickerMatch[1],
                        price: parseFloat(priceMatch[1]),
                        time: timeMatch[1]
                    };
                } else {
                    throw new Error('Invalid data format');
                }
            } catch (e) {
                error = 'Error parsing data: ' + e.message;
            }
        };

        socket.onclose = () => {
            connected = false;
            console.log('WebSocket connection closed');
        };
    }

    function sendSymbol() {
        if (!connected) {
            error = 'Not connected to server';
            return;
        }
        if (!symbol) {
            error = 'Please enter a symbol';
            return;
        }
        socket.send(symbol.toUpperCase());
        error = null;
    }

    onMount(() => {
        connectWebSocket();
    });

    onDestroy(() => {
        if (socket) {
            socket.close();
        }
    });
</script>

<div class="container">
    <div class="status-indicator" class:connected>
        {connected ? '● Connected' : '○ Disconnected'}
    </div>

    <div class="input-area">
        <input 
            type="text" 
            bind:value={symbol}
            placeholder="Enter stock symbol (e.g., MSFT)"
            on:keydown={(e) => e.key === 'Enter' && sendSymbol()}
        />
        <button on:click={sendSymbol} disabled={!connected}>
            Get Stock Data
        </button>
    </div>

    {#if error}
        <div class="error">{error}</div>
    {/if}

    {#if stockData}
        <div class="stock-info">
            <div class="data-row">
                <span class="label">Ticker:</span>
                <span class="value">{stockData.ticker}</span>
            </div>
            <div class="data-row">
                <span class="label">Price:</span>
                <span class="value">${stockData.price.toFixed(2)}</span>
            </div>
            <div class="data-row">
                <span class="label">Time:</span>
                <span class="value">{stockData.time}</span>
            </div>
        </div>
    {/if}
</div>

<style>
    .container {
        max-width: 600px;
        margin: 2rem auto;
        padding: 1rem;
        font-family: Arial, sans-serif;
    }

    .status-indicator {
        margin-bottom: 1rem;
        padding: 0.5rem;
        border-radius: 4px;
        text-align: center;
        color: #dc3545;
    }

    .status-indicator.connected {
        color: #28a745;
    }

    .input-area {
        display: flex;
        gap: 0.5rem;
        margin-bottom: 1rem;
    }

    input {
        flex: 1;
        padding: 0.5rem;
        border: 1px solid #ccc;
        border-radius: 4px;
    }

    button {
        padding: 0.5rem 1rem;
        background-color: #007bff;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
    }

    button:disabled {
        background-color: #ccc;
        cursor: not-allowed;
    }

    .error {
        color: #dc3545;
        padding: 0.5rem;
        margin: 1rem 0;
        background-color: #f8d7da;
        border-radius: 4px;
    }

    .stock-info {
        margin-top: 1rem;
        padding: 1rem;
        background-color: #f8f9fa;
        border-radius: 4px;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }

    .data-row {
        display: flex;
        justify-content: space-between;
        padding: 0.5rem 0;
        border-bottom: 1px solid #eee;
    }

    .data-row:last-child {
        border-bottom: none;
    }

    .label {
        font-weight: bold;
        color: #495057;
    }

    .value {
        color: #007bff;
    }
</style>