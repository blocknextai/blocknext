package server

import (
	"context"
	"time"

	"github.com/google/uuid"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	triggersContract "github.com/blocknextai/platform-api/internal/modules/triggers/contract"
)

type triggerOutput struct {
	ID               string    `json:"id" jsonschema:"trigger id"`
	Type             string    `json:"type" jsonschema:"trigger type: schedule or webhook"`
	ExecutionContext string    `json:"executionContext" jsonschema:"what the trigger runs (workflow, library item, ...)"`
	ContextItemID    string    `json:"contextItemId" jsonschema:"id of the item the trigger runs"`
	CronPattern      string    `json:"cronPattern,omitempty" jsonschema:"cron pattern of a schedule trigger"`
	Timezone         string    `json:"timezone,omitempty" jsonschema:"timezone of a schedule trigger"`
	IsActive         bool      `json:"isActive" jsonschema:"whether the trigger starts runs"`
	UpdatedAt        time.Time `json:"updatedAt" jsonschema:"last update time"`
}

type listTriggersOutput struct {
	Total    int64           `json:"total" jsonschema:"total number of triggers"`
	Triggers []triggerOutput `json:"triggers" jsonschema:"the triggers in this page"`
}

type setTriggerActiveInput struct {
	TriggerID string `json:"triggerId" jsonschema:"the trigger id"`
	IsActive  bool   `json:"isActive" jsonschema:"true to start runs from this trigger, false to stop them"`
}

type setTriggerActiveOutput struct {
	Trigger triggerOutput `json:"trigger" jsonschema:"the updated trigger"`
}

func (p *serverProvider) registerTriggerTools(server *mcpsdk.Server) {
	addTool(p, server, &mcpsdk.Tool{
		Name:        serverID + "_list_triggers",
		Title:       "List triggers",
		Description: "List the schedule and webhook triggers of the organization the access token acts on.",
		Annotations: readOnly("List triggers"),
	}, platformReadScope, p.listTriggers)

	addTool(p, server, &mcpsdk.Tool{
		Name:        serverID + "_set_trigger_active",
		Title:       "Activate or deactivate a trigger",
		Description: "Turn a trigger on or off. An inactive trigger no longer starts runs of the item it points to.",
		Annotations: mutating("Activate or deactivate a trigger"),
	}, platformWriteScope, p.setTriggerActive)
}

func (p *serverProvider) listTriggers(ctx context.Context, req *mcpsdk.CallToolRequest, input listInput) (*mcpsdk.CallToolResult, listTriggersOutput, error) {
	authenticated, err := requireCaller(req, platformReadScope)
	if err != nil {
		return nil, listTriggersOutput{}, err
	}

	offset, limit := normalizePagination(input.Offset, input.Limit)

	triggers, total, err := p.deps.TriggerService.GetAllByOrganizationID(ctx, authenticated.OrganizationID, input.Search, offset, limit)
	if err != nil {
		return nil, listTriggersOutput{}, err
	}

	outputs := make([]triggerOutput, 0, len(triggers))
	for _, trigger := range triggers {
		outputs = append(outputs, toTriggerOutput(trigger))
	}

	return nil, listTriggersOutput{
		Total:    total,
		Triggers: outputs,
	}, nil
}

func (p *serverProvider) setTriggerActive(ctx context.Context, req *mcpsdk.CallToolRequest, input setTriggerActiveInput) (*mcpsdk.CallToolResult, setTriggerActiveOutput, error) {
	authenticated, err := requireCaller(req, platformWriteScope)
	if err != nil {
		return nil, setTriggerActiveOutput{}, err
	}

	triggerID, err := uuid.Parse(input.TriggerID)
	if err != nil {
		return nil, setTriggerActiveOutput{}, ErrInvalidTriggerID
	}

	trigger, err := p.deps.TriggerService.SetActive(ctx, authenticated.OrganizationID, triggerID, input.IsActive)
	if err != nil {
		return nil, setTriggerActiveOutput{}, err
	}

	return nil, setTriggerActiveOutput{
		Trigger: toTriggerOutput(trigger),
	}, nil
}

func toTriggerOutput(trigger *triggersContract.Trigger) triggerOutput {
	return triggerOutput{
		ID:               trigger.ID.String(),
		Type:             trigger.Type.String(),
		ExecutionContext: string(trigger.ExecutionContext),
		ContextItemID:    trigger.ContextItemID.String(),
		CronPattern:      stringValue(trigger.CronPattern),
		Timezone:         stringValue(trigger.Timezone),
		IsActive:         trigger.IsActive,
		UpdatedAt:        trigger.UpdatedAt,
	}
}
