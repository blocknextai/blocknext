package getallwebhooksources

import (
	"context"

	nodeEngineApplicationAdapters "github.com/blocknextai/platform-api/internal/nodeengine/application/adapters"
)

type Handler struct {
	adapterService     nodeEngineApplicationAdapters.AdapterService
	webhookURLTemplate string
}

func New(
	adapterService nodeEngineApplicationAdapters.AdapterService,
	webhookURLTemplate string,
) *Handler {
	return &Handler{
		adapterService:     adapterService,
		webhookURLTemplate: webhookURLTemplate,
	}
}

func (h *Handler) Handle(_ context.Context, _ *GetAllWebhookSourcesQuery) (*GetAllWebhookSourcesResponse, error) {
	adapters := h.adapterService.GetAllAdapters()
	return MapGetAllWebhookSourcesQueryToGetAllWebhookSourcesResponse(adapters, h.webhookURLTemplate), nil
}
