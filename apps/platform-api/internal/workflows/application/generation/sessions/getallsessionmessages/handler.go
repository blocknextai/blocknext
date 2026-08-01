package getallsessionmessages

import (
	"context"

	generationDomainMessages "github.com/blocknextai/platform-api/internal/workflows/domain/generation/messages"
	generationDomainSessions "github.com/blocknextai/platform-api/internal/workflows/domain/generation/sessions"
)

type Handler struct {
	sessionRepository generationDomainSessions.SessionRepository
	messageRepository generationDomainMessages.MessageRepository
}

func New(
	sessionRepository generationDomainSessions.SessionRepository,
	messageRepository generationDomainMessages.MessageRepository,
) *Handler {
	return &Handler{
		sessionRepository: sessionRepository,
		messageRepository: messageRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetAllSessionMessagesQuery) (*GetAllSessionMessagesResponse, error) {
	session, err := h.sessionRepository.GetByIDAndOrganizationID(ctx, request.SessionID, request.OrganizationID)
	if err != nil {
		return nil, err
	}

	messages, totalCount, err := h.messageRepository.GetAllBySessionID(
		ctx,
		session.ID,
		request.Pagination.Offset,
		request.Pagination.Limit,
	)
	if err != nil {
		return nil, err
	}

	return &GetAllSessionMessagesResponse{
		Items:      MapGetAllSessionMessagesQueryToGetAllSessionMessagesResponse(messages),
		TotalCount: totalCount,
	}, nil
}
