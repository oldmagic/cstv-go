package tests

import (
	"log"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketLoad(t *testing.T) {
	serverURL := "ws://localhost:8080/ws"
	numClients := 1000
	var wg sync.WaitGroup

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			u, _ := url.Parse(serverURL)
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			assert.NoError(t, err)
			defer conn.Close()

			msg := "Hello from client " + strconv.Itoa(id)
			err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
			assert.NoError(t, err)

			_, response, err := conn.ReadMessage()
			assert.NoError(t, err)
			log.Printf("Client %d received: %s", id, string(response))
		}(i)
	}

	wg.Wait()
}
