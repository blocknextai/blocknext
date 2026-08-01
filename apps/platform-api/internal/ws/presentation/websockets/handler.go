package websockets

import (
	"log/slog"
	"sync"
	"time"

	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	wsApplicationHub "github.com/blocknextai/platform-api/internal/ws/application/hub"
	wsDomainConnections "github.com/blocknextai/platform-api/internal/ws/domain/connections"
	"github.com/gofiber/contrib/v3/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = 54 * time.Second
)

func NewWebSocketHandler(hub wsApplicationHub.HubService) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		userID := commonHTTP.GetUserIDFromWebSocket(c)
		organizationID := commonHTTP.GetOrganizationIDFromWebSocket(c)

		conn := wsDomainConnections.New(userID, organizationID)
		hub.Register(conn)

		var closeOnce sync.Once
		closeWs := func() {
			closeOnce.Do(func() {
				if err := c.Close(); err != nil {
					slog.Warn("websocket close failed",
						"connectionId", conn.ID,
						"error", err)
				}
			})
		}

		defer func() {
			hub.Unregister(conn)
			closeWs()
		}()

		go writePump(c, conn, closeWs)

		if err := c.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			slog.Warn("websocket set read deadline failed",
				"connectionId", conn.ID,
				"error", err)
		}
		c.SetPongHandler(func(string) error {
			if err := c.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
				slog.Warn("websocket set read deadline failed",
					"connectionId", conn.ID,
					"error", err)
			}
			return nil
		})

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					slog.Error("websocket read error",
						"connectionId", conn.ID,
						"error", err,
					)
				}
				break
			}
		}
	}
}

func writePump(c *websocket.Conn, conn *wsDomainConnections.Connection, closeWs func()) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		closeWs()
	}()

	for {
		select {
		case message, ok := <-conn.Send:
			if err := c.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				slog.Warn("websocket set write deadline failed",
					"connectionId", conn.ID,
					"error", err)
			}
			if !ok {
				if err := c.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					slog.Warn("websocket close message failed",
						"connectionId", conn.ID,
						"error", err)
				}
				return
			}
			if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				slog.Warn("websocket set write deadline failed",
					"connectionId", conn.ID,
					"error", err)
			}
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
