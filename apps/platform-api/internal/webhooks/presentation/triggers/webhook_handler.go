package triggers

import (
	"github.com/blocknextai/go-packages/json"
	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	nodeEngineDomainAdapters "github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
	taskRunnerApplicationWebhooks "github.com/blocknextai/platform-api/internal/taskrunner/application/webhooks"
	triggersApplicationWebhooks "github.com/blocknextai/platform-api/internal/triggers/application/webhooks"
	"github.com/gofiber/fiber/v3"
)

type WebhookRequest struct {
	Source string `uri:"source"`
	Token  string `uri:"token"`
}

func NewWebhookHandler(processor taskRunnerApplicationWebhooks.WebhookProcessor) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(WebhookRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		var payload map[string]any
		if err := json.Unmarshal(c.Body(), &payload); err != nil {
			payload = make(map[string]any)
		}

		result, err := processor.ProcessWebhook(c.RequestCtx(), &triggersApplicationWebhooks.Request{
			WebhookToken: request.Token,
			Source:       request.Source,
			Payload:      payload,
			VerificationRequest: &nodeEngineDomainAdapters.VerificationRequest{
				Method:      c.Method(),
				Headers:     c.GetReqHeaders(),
				Body:        c.Body(),
				QueryParams: c.Queries(),
			},
		})
		if err != nil {
			return err
		}

		if result.Verification != nil {
			return c.Status(result.Verification.StatusCode).Send(result.Verification.Body)
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("webhook trigger processed")))
	}
}
