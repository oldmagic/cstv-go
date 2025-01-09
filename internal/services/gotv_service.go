package services

import (
    "sync"

    "github.com/sirupsen/logrus"
)

// GOTVService manages active WebSocket clients and relays messages.
type GOTVService struct {
    mu      sync.Mutex
    clients map[string]chan string
}

// NewGOTVService initializes a new GOTVService instance.
func NewGOTVService() *GOTVService {
    return &GOTVService{
        clients: make(map[string]chan string),
    }
}

// RegisterClient adds a new WebSocket client.
func (s *GOTVService) RegisterClient(id string, messageChannel chan string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    s.clients[id] = messageChannel
    logrus.Infof("Client %s registered", id)
}

// UnregisterClient removes a WebSocket client.
func (s *GOTVService) UnregisterClient(id string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    delete(s.clients, id)
    logrus.Infof("Client %s unregistered", id)
}

// BroadcastMessage sends a message to all registered clients.
func (s *GOTVService) BroadcastMessage(message string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    for id, ch := range s.clients {
        select {
        case ch <- message:
            logrus.Infof("Message sent to client %s", id)
        default:
            logrus.Warnf("Client %s channel full, dropping message", id)
        }
    }
}
