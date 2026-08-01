package generation

import (
	"bufio"
	"log/slog"
	"strings"

	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/llm/streamingchat"
	"github.com/blocknextai/platform-api/internal/workflows/application/generation/chat"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

type SendMessageRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	SessionID      uuid.UUID `uri:"sessionId"`
	Message        string    `json:"message"`
}

func NewSendMessageHandler(chatService chat.ChatService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if chatService == nil {
			return chat.ErrProviderNotReady
		}

		request := new(SendMessageRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		streamCh, err := chatService.SendMessage(c.RequestCtx(), request.OrganizationID, request.SessionID, request.Message)
		if err != nil {
			return err
		}

		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("X-Accel-Buffering", "no")

		c.RequestCtx().SetBodyStreamWriter(func(w *bufio.Writer) {
			write := func(s string) {
				if _, err := w.WriteString(s); err != nil {
					slog.WarnContext(c.RequestCtx(), "sse write failed",
						"component", "SendMessageStream",
						"error", err)
				}
			}
			for chunk := range streamCh {
				switch chunk.Type {
				case streamingchat.ChunkTypeText:
					sseData := strings.ReplaceAll(chunk.Content, "\n", "\ndata: ")
					write("event: message\ndata: ")
					write(sseData)
					write("\n\n")
				case streamingchat.ChunkTypeDone:
					write("event: done\ndata: [DONE]\n\n")
				case streamingchat.ChunkTypeError:
					errMsg := "unknown error"
					if chunk.Error != nil {
						errMsg = chunk.Error.Error()
					}
					write("event: error\ndata: ")
					write(errMsg)
					write("\n\n")
				}
				if err := w.Flush(); err != nil {
					return
				}
			}
		})

		return nil
	}
}

var _ fasthttp.StreamWriter = func(w *bufio.Writer) {}
