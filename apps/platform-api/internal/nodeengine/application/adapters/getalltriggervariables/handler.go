package getalltriggervariables

import (
	"context"

	nodeEngineApplicationAdapters "github.com/blocknextai/platform-api/internal/nodeengine/application/adapters"
)

type Handler struct {
	adapterService nodeEngineApplicationAdapters.AdapterService
}

func New(
	adapterService nodeEngineApplicationAdapters.AdapterService,
) *Handler {
	return &Handler{
		adapterService: adapterService,
	}
}

func (h *Handler) Handle(_ context.Context, _ *GetAllTriggerVariablesQuery) (*GetAllTriggerVariablesResponse, error) {
	triggerVariables := h.adapterService.GetTriggerVariables()
	return MapGetAllTriggerVariablesQueryToGetAllTriggerVariablesResponse(triggerVariables), nil
}
