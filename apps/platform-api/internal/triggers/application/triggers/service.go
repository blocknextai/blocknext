package application

import (
	"context"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
	"github.com/google/uuid"
)

type TriggerService interface {
	GetAllActive(
		ctx context.Context,
	) ([]*triggersDomainTriggers.Trigger, error)

	Create(
		ctx context.Context,
		organizationID uuid.UUID,
		triggeredByUserID *uuid.UUID,
		executionContext commonDomain.ExecutionContext,
		contextItemID uuid.UUID,
		triggerType triggersDomainTriggers.TriggerType,
		cronPattern *string,
		timezone *string,
		runtimeConfig *triggersDomainTriggers.RuntimeConfig,
	) (*triggersDomainTriggers.Trigger, string, error)
}

type triggerService struct {
	triggerRepository triggersDomainTriggers.TriggerRepository
}

func NewTriggerService(
	triggerRepository triggersDomainTriggers.TriggerRepository,
) TriggerService {
	return &triggerService{
		triggerRepository: triggerRepository,
	}
}

func (s *triggerService) GetAllActive(
	ctx context.Context,
) ([]*triggersDomainTriggers.Trigger, error) {
	return s.triggerRepository.GetAllActive(ctx)
}

func (s *triggerService) Create(
	ctx context.Context,
	organizationID uuid.UUID,
	triggeredByUserID *uuid.UUID,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	triggerType triggersDomainTriggers.TriggerType,
	cronPattern *string,
	timezone *string,
	runtimeConfig *triggersDomainTriggers.RuntimeConfig,
) (*triggersDomainTriggers.Trigger, string, error) {
	var tokenHash *string
	var plainToken string
	if triggerType == triggersDomainTriggers.TriggerTypeWebhook {
		plain, hash, err := GenerateWebhookToken()
		if err != nil {
			return nil, "", err
		}
		plainToken = plain
		tokenHash = new(hash)
	}

	trigger, err := triggersDomainTriggers.New(
		organizationID,
		triggeredByUserID,
		executionContext,
		contextItemID,
		triggerType,
		cronPattern,
		timezone,
		tokenHash,
		nil,
		runtimeConfig,
		true,
	)
	if err != nil {
		return nil, "", err
	}

	err = s.triggerRepository.Create(ctx, trigger)
	if err != nil {
		return nil, "", err
	}

	return trigger, plainToken, nil
}
