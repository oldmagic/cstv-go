package handlers

import (
    "log"

    "github.com/oldmagic/cstv-go/internal/services"
    "github.com/gofiber/websocket/v2"
    "github.com/google/uuid"
)

// WebSocketHandler handles WebSocket connections with GOTVService.
func WebSocketHandler(service *services.GOTVService) func(*websocket.Conn) {
    return func(c *websocket.Conn) {
        defer c.Close()

        // Generate a unique client ID
        clientID := uuid.New().String()
        messageChannel := make(chan string, 10)

        // Register client
        service.RegisterClient(clientID, messageChannel)
        defer service.UnregisterClient(clientID)

        // Goroutine to send messages to the client
        go func() {
            for msg := range messageChannel {
                if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
                    log.Println("Write error:", err)
                    break
                }
            }
        }()

        // Read loop
        for {
            _, msg, err := c.ReadMessage()
            if err != nil {
                log.Println("WebSocket error:", err)
                break
            }

            log.Println("Received message:", string(msg))
            service.BroadcastMessage(string(msg))
        }
    }
}
