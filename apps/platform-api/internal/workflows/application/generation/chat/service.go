package chat

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/platform-api/internal/llm/streamingchat"
	"github.com/blocknextai/platform-api/internal/workflows/application/generation/credentialschema"
	"github.com/blocknextai/platform-api/internal/workflows/application/generation/nodeschema"
	"github.com/blocknextai/platform-api/internal/workflows/application/generation/triggervariables"
	generationDomainMessages "github.com/blocknextai/platform-api/internal/workflows/domain/generation/messages"
	generationDomainSessions "github.com/blocknextai/platform-api/internal/workflows/domain/generation/sessions"
	"github.com/google/uuid"
)

type ChatService interface {
	SendMessage(ctx context.Context, organizationID uuid.UUID, sessionID uuid.UUID, userMessage string) (<-chan streamingchat.Chunk, error)
}

type chatService struct {
	systemInstruction              string
	sessionRepository              generationDomainSessions.SessionRepository
	messageRepository              generationDomainMessages.MessageRepository
	aiProvider                     streamingchat.Provider
	nodeSchemaContextBuilder       *nodeschema.NodeSchemaContextBuilder
	credentialSchemaContextBuilder *credentialschema.CredentialSchemaContextBuilder
	triggerVariablesContextBuilder *triggervariables.TriggerVariablesContextBuilder
}

func NewChatService(
	systemInstruction string,
	sessionRepository generationDomainSessions.SessionRepository,
	messageRepository generationDomainMessages.MessageRepository,
	aiProvider streamingchat.Provider,
	nodeSchemaContextBuilder *nodeschema.NodeSchemaContextBuilder,
	credentialSchemaContextBuilder *credentialschema.CredentialSchemaContextBuilder,
	triggerVariablesContextBuilder *triggervariables.TriggerVariablesContextBuilder,
) ChatService {
	return &chatService{
		systemInstruction:              systemInstruction,
		sessionRepository:              sessionRepository,
		messageRepository:              messageRepository,
		aiProvider:                     aiProvider,
		nodeSchemaContextBuilder:       nodeSchemaContextBuilder,
		credentialSchemaContextBuilder: credentialSchemaContextBuilder,
		triggerVariablesContextBuilder: triggerVariablesContextBuilder,
	}
}

func (s *chatService) buildSystemInstruction() string {
	prompt := s.systemInstruction
	prompt = strings.ReplaceAll(prompt, "{AVAILABLE_NODES_JSON}", s.nodeSchemaContextBuilder.Build())
	prompt = strings.ReplaceAll(prompt, "{AVAILABLE_CREDENTIALS_JSON}", s.credentialSchemaContextBuilder.Build())
	prompt = strings.ReplaceAll(prompt, "{AVAILABLE_TRIGGER_VARIABLES_JSON}", s.triggerVariablesContextBuilder.Build())
	return prompt
}

var (
	ErrMessageIsRequired = apperror.Validation("message is required")
	ErrProviderNotReady  = apperror.Unavailable("provider not ready")
)

func (s *chatService) SendMessage(ctx context.Context, organizationID uuid.UUID, sessionID uuid.UUID, userMessage string) (<-chan streamingchat.Chunk, error) {
	if strings.TrimSpace(userMessage) == "" {
		return nil, ErrMessageIsRequired
	}

	session, err := s.sessionRepository.GetByIDAndOrganizationID(ctx, sessionID, organizationID)
	if err != nil {
		return nil, err
	}

	userMsg, err := generationDomainMessages.New(session.ID, "user", userMessage, nil)
	if err != nil {
		return nil, err
	}

	if err := s.messageRepository.Create(ctx, userMsg); err != nil {
		return nil, err
	}

	session.UpdatedAt = time.Now().UTC()
	if err := s.sessionRepository.Update(ctx, session); err != nil {
		slog.Error("Failed to update session updated_at",
			"component", "Generation",
			"sessionId", session.ID,
			"error", err,
		)
	}

	existingMessages, _, err := s.messageRepository.GetAllBySessionID(ctx, session.ID, 0, 1000)
	if err != nil {
		return nil, err
	}

	aiMessages := make([]streamingchat.Message, 0, len(existingMessages))
	for _, msg := range existingMessages {
		aiMessages = append(aiMessages, streamingchat.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	fullSystemInstruction := s.buildSystemInstruction()

	streamCh, err := s.aiProvider.StreamChat(ctx, fullSystemInstruction, aiMessages)
	if err != nil {
		return nil, err
	}

	outCh := make(chan streamingchat.Chunk, 64)

	go func() {
		defer close(outCh)

		var fullContent strings.Builder

		for chunk := range streamCh {
			switch chunk.Type {
			case streamingchat.ChunkTypeText:
				fullContent.WriteString(chunk.Content)
				select {
				case outCh <- chunk:
				case <-ctx.Done():
					return
				}

			case streamingchat.ChunkTypeDone:
				assistantMsg, err := generationDomainMessages.New(
					session.ID, "model", fullContent.String(), nil,
				)
				if err != nil {
					slog.Error("Failed to create assistant message entity",
						"component", "Generation",
						"sessionId", session.ID,
						"error", err,
					)
				} else {
					if err := s.messageRepository.Create(ctx, assistantMsg); err != nil {
						slog.Error("Failed to save assistant message",
							"component", "Generation",
							"sessionId", session.ID,
							"error", err,
						)
					}
				}

				select {
				case outCh <- chunk:
				case <-ctx.Done():
					return
				}

			case streamingchat.ChunkTypeError:
				select {
				case outCh <- chunk:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return outCh, nil
}
