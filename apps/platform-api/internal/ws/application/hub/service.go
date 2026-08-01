package hub

import (
	"context"
	"log/slog"
	"sync"

	"github.com/blocknextai/platform-api/internal/realtime"
	wsDomainConnections "github.com/blocknextai/platform-api/internal/ws/domain/connections"
	wsDomainRooms "github.com/blocknextai/platform-api/internal/ws/domain/rooms"
	"github.com/google/uuid"
)

type HubService interface {
	Register(conn *wsDomainConnections.Connection)
	Unregister(conn *wsDomainConnections.Connection)
	Shutdown()
}

type hubService struct {
	rooms                 map[uuid.UUID]*wsDomainRooms.Room
	mu                    sync.RWMutex
	broadcaster           realtime.Broadcaster
	maxConnectionsPerRoom int
	ctx                   context.Context
	cancel                context.CancelFunc
}

func NewHubService(broadcaster realtime.Broadcaster, maxConnectionsPerRoom int) HubService {
	ctx, cancel := context.WithCancel(context.Background())
	return &hubService{
		rooms:                 make(map[uuid.UUID]*wsDomainRooms.Room),
		broadcaster:           broadcaster,
		maxConnectionsPerRoom: maxConnectionsPerRoom,
		ctx:                   ctx,
		cancel:                cancel,
	}
}

func (h *hubService) Register(conn *wsDomainConnections.Connection) {
	room, isNew := h.addConnection(conn)
	if room == nil {
		return
	}

	if !isNew {
		return
	}

	ctx, cancel := context.WithCancel(h.ctx)
	room.SetCancel(cancel)

	if err := h.startSubscription(ctx, room); err != nil {
		h.mu.Lock()
		delete(h.rooms, conn.OrganizationID)
		h.mu.Unlock()
		room.CloseAllConnections()
	}
}

func (h *hubService) addConnection(conn *wsDomainConnections.Connection) (*wsDomainRooms.Room, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.rooms[conn.OrganizationID]
	if !exists {
		room = wsDomainRooms.New(conn.OrganizationID, h.maxConnectionsPerRoom)
		h.rooms[conn.OrganizationID] = room
		room.AddConnection(conn)
		return room, true
	}

	if !room.AddConnection(conn) {
		slog.Warn("max connections per room reached, rejecting connection",
			"organizationId", conn.OrganizationID,
			"connectionId", conn.ID,
			"maxConnectionsPerRoom", h.maxConnectionsPerRoom,
		)
		conn.Close()
		return nil, false
	}

	return room, false
}

func (h *hubService) Unregister(conn *wsDomainConnections.Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.rooms[conn.OrganizationID]
	if !exists {
		return
	}

	remaining := room.RemoveConnection(conn.ID)
	conn.Close()

	if remaining == 0 {
		room.Stop()
		delete(h.rooms, conn.OrganizationID)
	}
}

func (h *hubService) startSubscription(ctx context.Context, room *wsDomainRooms.Room) error {
	ch, err := h.broadcaster.Subscribe(ctx, room.OrganizationID)
	if err != nil {
		slog.Error("failed to subscribe to organization channel",
			"organizationId", room.OrganizationID,
			"error", err,
		)
		room.Stop()
		return err
	}

	go func() {
		for msg := range ch {
			room.Broadcast([]byte(msg))
		}
	}()

	return nil
}

func (h *hubService) Shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.cancel()

	for _, room := range h.rooms {
		room.Stop()
		room.CloseAllConnections()
	}
	h.rooms = make(map[uuid.UUID]*wsDomainRooms.Room)
}
