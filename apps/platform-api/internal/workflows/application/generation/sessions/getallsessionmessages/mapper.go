package getallsessionmessages

import (
	generationDomainMessages "github.com/blocknextai/platform-api/internal/workflows/domain/generation/messages"
)

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
