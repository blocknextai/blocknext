package generation

import (
	"context"
	"time"

	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	generationDomainMessages "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/generation/messages"
)

type GetAllSessionMessagesQuery struct {
	OrganizationID uuid.UUID
	SessionID      uuid.UUID
	Pagination     resultPkg.PaginationRequest
}

type MessageResponse struct {
	ID        uuid.UUID      `json:"id"`
	SessionID uuid.UUID      `json:"sessionId"`
	Role      string         `json:"role"`
	Content   string         `json:"content"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
}

type GetAllSessionMessagesResponse struct {
	Items      []*MessageResponse
	TotalCount int64
}

func (s *Service) GetAllSessionMessages(ctx context.Context, request *GetAllSessionMessagesQuery) (*GetAllSessionMessagesResponse, error) {
	session, err := s.sessionRepository.GetByIDAndOrganizationID(ctx, request.SessionID, request.OrganizationID)
	if err != nil {
		return nil, err
	}

	messages, totalCount, err := s.messageRepository.GetAllBySessionID(
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

func MapGetAllSessionMessagesQueryToGetAllSessionMessagesResponse(messages []*generationDomainMessages.GenerationMessage) []*MessageResponse {
	response := make([]*MessageResponse, 0, len(messages))
	for _, msg := range messages {
		response = append(response, &MessageResponse{
			ID:        msg.ID,
			SessionID: msg.SessionID,
			Role:      msg.Role,
			Content:   msg.Content,
			Metadata:  msg.Metadata,
			CreatedAt: msg.CreatedAt,
		})
	}

	return response
}
