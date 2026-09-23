package application

import (
	"context"

	"github.com/google/uuid"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/domain/triggers"
)

type TriggerService interface {
	GetAllActive(
		ctx context.Context,
	) ([]*triggersDomainTriggers.Trigger, error)

	GetAllByOrganizationID(
		ctx context.Context,
		organizationID uuid.UUID,
		searchQuery string,
		offset int,
		limit int,
	) ([]*triggersDomainTriggers.Trigger, int64, error)

	SetActive(
		ctx context.Context,
		organizationID uuid.UUID,
		triggerID uuid.UUID,
		isActive bool,
	) (*triggersDomainTriggers.Trigger, error)

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

func (s *triggerService) GetAllByOrganizationID(
	ctx context.Context,
	organizationID uuid.UUID,
	searchQuery string,
	offset int,
	limit int,
) ([]*triggersDomainTriggers.Trigger, int64, error) {
	return s.triggerRepository.GetAllByOrganizationID(ctx, organizationID, searchQuery, offset, limit)
}

func (s *triggerService) SetActive(
	ctx context.Context,
	organizationID uuid.UUID,
	triggerID uuid.UUID,
	isActive bool,
) (*triggersDomainTriggers.Trigger, error) {
	trigger, err := s.triggerRepository.GetByIDAndOrganizationID(ctx, triggerID, organizationID)
	if err != nil {
		return nil, err
	}

	var updated *triggersDomainTriggers.Trigger
	switch trigger.Type {
	case triggersDomainTriggers.TriggerTypeSchedule:
		updated, err = trigger.UpdateSchedule(isActive, trigger.CronPattern, trigger.Timezone, trigger.RuntimeConfig)
	case triggersDomainTriggers.TriggerTypeWebhook:
		updated, err = trigger.UpdateWebhook(isActive, trigger.WebhookTokenHash, trigger.WebhookSecret, trigger.RuntimeConfig)
	default:
		return nil, triggersDomainTriggers.ErrUnsupportedType
	}
	if err != nil {
		return nil, err
	}

	if err := s.triggerRepository.Update(ctx, updated); err != nil {
		return nil, err
	}

	return updated, nil
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
