package adapters

import (
	"context"
	"strings"

	adaptersDomain "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/adapters"
)

type GetAllWebhookSourcesQuery struct{}

type WebhookSourceResponse struct {
	Source               string `json:"source"`
	Name                 string `json:"name"`
	SupportsVerification bool   `json:"supportsVerification"`
	WebhookURL           string `json:"webhookUrl"`
}

type GetAllWebhookSourcesResponse = []WebhookSourceResponse

func (s *Service) GetAllWebhookSources(ctx context.Context, _ *GetAllWebhookSourcesQuery) (*GetAllWebhookSourcesResponse, error) {
	adapters := s.adapterService.GetAllAdapters()
	return MapGetAllWebhookSourcesQueryToGetAllWebhookSourcesResponse(adapters, s.webhookURLTemplate), nil
}

func MapGetAllWebhookSourcesQueryToGetAllWebhookSourcesResponse(
	adapters []adaptersDomain.TriggerAdapter,
	webhookURLTemplate string,
) *GetAllWebhookSourcesResponse {
	response := make(GetAllWebhookSourcesResponse, 0, len(adapters))
	for _, adapter := range adapters {
		url := strings.ReplaceAll(webhookURLTemplate, "{source}", adapter.GetID())
		_, supportsVerification := adaptersDomain.AsWebhookVerifier(adapter)
		response = append(response, WebhookSourceResponse{
			Source:               adapter.GetID(),
			Name:                 adapter.GetName(),
			SupportsVerification: supportsVerification,
			WebhookURL:           url,
		})
	}
	return &response
}
