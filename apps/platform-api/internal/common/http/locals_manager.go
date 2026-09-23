package http

import (
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const (
	userIDContextKey         = "userId"
	sessionIDContextKey      = "sessionId"
	organizationIDContextKey = "organizationId"
)

func SetUserID(c fiber.Ctx, userID uuid.UUID) {
	c.Locals(userIDContextKey, userID)
}

func GetUserID(c fiber.Ctx) uuid.UUID {
	if value := c.Locals(userIDContextKey); value != nil {
		if userID, ok := value.(uuid.UUID); ok {
			return userID
		}
	}
	return uuid.Nil
}

func SetSessionID(c fiber.Ctx, sessionID uuid.UUID) {
	c.Locals(sessionIDContextKey, sessionID)
}

func GetSessionID(c fiber.Ctx) uuid.UUID {
	if value := c.Locals(sessionIDContextKey); value != nil {
		if sessionID, ok := value.(uuid.UUID); ok {
			return sessionID
		}
	}
	return uuid.Nil
}

func SetOrganizationID(c fiber.Ctx, organizationID uuid.UUID) {
	c.Locals(organizationIDContextKey, organizationID)
}

func GetOrganizationID(c fiber.Ctx) uuid.UUID {
	if value := c.Locals(organizationIDContextKey); value != nil {
		if organizationID, ok := value.(uuid.UUID); ok {
			return organizationID
		}
	}
	return uuid.Nil
}

func GetUserIDFromWebSocket(c *websocket.Conn) uuid.UUID {
	if value := c.Locals(userIDContextKey); value != nil {
		if userID, ok := value.(uuid.UUID); ok {
			return userID
		}
	}
	return uuid.Nil
}

func GetOrganizationIDFromWebSocket(c *websocket.Conn) uuid.UUID {
	if value := c.Locals(organizationIDContextKey); value != nil {
		if organizationID, ok := value.(uuid.UUID); ok {
			return organizationID
		}
	}
	return uuid.Nil
}
