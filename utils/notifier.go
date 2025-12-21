package utils

import (
	"sync"
)

// Simple in-memory notifier for server-sent events
type Notifier struct {
	clients map[chan string]struct{}
	mu      sync.Mutex
}

func NewNotifier() *Notifier {
	return &Notifier{clients: make(map[chan string]struct{})}
}

var NotifierInstance = NewNotifier()

func (n *Notifier) AddClient(ch chan string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.clients[ch] = struct{}{}
}

func (n *Notifier) RemoveClient(ch chan string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.clients, ch)
	close(ch)
}

func (n *Notifier) Notify(msg string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for ch := range n.clients {
		// non-blocking send
		select {
		case ch <- msg:
		default:
		}
	}
}
