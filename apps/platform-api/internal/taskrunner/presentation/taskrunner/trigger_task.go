package taskrunner

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/dag"
	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/taskrunner/application/triggertask"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
)

type TriggerTaskRequest struct {
	OrganizationID   uuid.UUID                            `uri:"organizationId"`
	ExecutionContext commonDomain.ExecutionContext        `json:"executionContext"`
	ContextItemID    uuid.UUID                            `json:"contextItemId"`
	TriggerType      taskRunnerDomainTask.TaskTriggerType `json:"triggerType"`
	CronPattern      *string                              `json:"cronPattern"`

	RuntimePrompt string `json:"runtimePrompt"`

	Nodes []dag.Node `json:"nodes"`
}

func NewTriggerTaskHandler(handler *triggertask.TriggerTaskHandler) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(TriggerTaskRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &triggertask.TriggerTaskCommand{
			TriggeredByUserID: new(userID),
			OrganizationID:    request.OrganizationID,
			ExecutionContext:  request.ExecutionContext,
			ContextItemID:     request.ContextItemID,
			TriggerType:       request.TriggerType,
			CronPattern:       request.CronPattern,
			RuntimePrompt:     request.RuntimePrompt,
			Nodes:             request.Nodes,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("your flow has started running")))
	}
}

type TriggerTaskWithAPIWorkflowRequest struct {
	WorkflowID uuid.UUID `uri:"workflowId"`

	RuntimePrompt string `json:"runtimePrompt"`

	Nodes []dag.Node `json:"nodes"`
}

func NewTriggerTaskWithAPIWorkflowHandler(handler *triggertask.TriggerTaskHandler) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(TriggerTaskWithAPIWorkflowRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		organizationID := commonHTTP.GetOrganizationID(c)
		if organizationID == uuid.Nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &triggertask.TriggerTaskCommand{
			TriggeredByUserID: nil,
			OrganizationID:    organizationID,
			ExecutionContext:  commonDomain.ExecutionContextWorkflow,
			ContextItemID:     request.WorkflowID,
			TriggerType:       taskRunnerDomainTask.TaskTriggerTypeAPI,
			RuntimePrompt:     request.RuntimePrompt,
			Nodes:             request.Nodes,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("api trigger processed")))
	}
}
