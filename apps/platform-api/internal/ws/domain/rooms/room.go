package rooms

import (
	"context"
	"log/slog"
	"sync"

	wsDomainConnections "github.com/blocknextai/platform-api/internal/ws/domain/connections"
	"github.com/google/uuid"
)

type Room struct {
	OrganizationID uuid.UUID
	connections    map[string]*wsDomainConnections.Connection
	mu             sync.RWMutex
	cancel         context.CancelFunc
	maxConnections int
}

func New(organizationID uuid.UUID, maxConnections int) *Room {
	return &Room{
		OrganizationID: organizationID,
		connections:    make(map[string]*wsDomainConnections.Connection),
		maxConnections: maxConnections,
	}
}

func (r *Room) AddConnection(conn *wsDomainConnections.Connection) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.maxConnections > 0 && len(r.connections) >= r.maxConnections {
		return false
	}
	r.connections[conn.ID] = conn
	return true
}

func (r *Room) RemoveConnection(connID string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.connections, connID)
	return len(r.connections)
}

func (r *Room) Broadcast(message []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, conn := range r.connections {
		select {
		case conn.Send <- message:
		default:
			slog.Warn("message dropped, send buffer full",
				"connectionId", conn.ID,
				"organizationId", r.OrganizationID,
			)
		}
	}
}

func (r *Room) SetCancel(cancel context.CancelFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cancel = cancel
}

func (r *Room) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cancel != nil {
		r.cancel()
		r.cancel = nil
	}
}

func (r *Room) CloseAllConnections() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, conn := range r.connections {
		conn.Close()
	}
}
