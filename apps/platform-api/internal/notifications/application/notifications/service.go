package notifications

import (
	"context"

	notificationsDomainRecipients "github.com/blocknextai/platform-api/internal/notifications/domain/notificationrecipients"
	notificationsDomainNotifications "github.com/blocknextai/platform-api/internal/notifications/domain/notifications"
	organizationsApplicationOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	"github.com/google/uuid"
)

type CreateInput struct {
	Type         string
	Level        notificationsDomainNotifications.Level
	AudienceType notificationsDomainNotifications.AudienceType
	AudienceID   uuid.UUID
	Title        string
	Body         *string
	Data         map[string]any
	ActionURL    *string
}

type NotificationService interface {
	Create(ctx context.Context, input CreateInput) error
}

type notificationService struct {
	notificationRepository          notificationsDomainNotifications.NotificationRepository
	notificationRecipientRepository notificationsDomainRecipients.NotificationRecipientRepository
	organizationUserService         organizationsApplicationOrganizationUsers.OrganizationUserService
}

func NewNotificationService(
	notificationRepository notificationsDomainNotifications.NotificationRepository,
	notificationRecipientRepository notificationsDomainRecipients.NotificationRecipientRepository,
	organizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService,
) NotificationService {
	return &notificationService{
		notificationRepository:          notificationRepository,
		notificationRecipientRepository: notificationRecipientRepository,
		organizationUserService:         organizationUserService,
	}
}

func (s *notificationService) Create(ctx context.Context, input CreateInput) error {
	notification, err := notificationsDomainNotifications.New(
		input.Type,
		input.Level,
		input.AudienceType,
		input.AudienceID,
		input.Title,
		input.Body,
		input.Data,
		input.ActionURL,
	)
	if err != nil {
		return err
	}

	if err := s.notificationRepository.Create(ctx, notification); err != nil {
		return err
	}

	userIDs, err := s.resolveRecipientUserIDs(ctx, notification.AudienceType, notification.AudienceID)
	if err != nil {
		return err
	}

	if len(userIDs) == 0 {
		return nil
	}

	recipients := make([]*notificationsDomainRecipients.NotificationRecipient, 0, len(userIDs))
	for _, userID := range userIDs {
		recipient, err := notificationsDomainRecipients.New(notification.ID, userID)
		if err != nil {
			return err
		}
		recipients = append(recipients, recipient)
	}

	return s.notificationRecipientRepository.CreateMany(ctx, recipients)
}

func (s *notificationService) resolveRecipientUserIDs(ctx context.Context, audienceType notificationsDomainNotifications.AudienceType, audienceID uuid.UUID) ([]uuid.UUID, error) {
	switch audienceType {
	case notificationsDomainNotifications.AudienceTypeUser:
		return []uuid.UUID{
			audienceID,
		}, nil
	case notificationsDomainNotifications.AudienceTypeOrganization:
		members, err := s.organizationUserService.GetAllByOrganizationID(ctx, audienceID)
		if err != nil {
			return nil, err
		}

		userIDs := make([]uuid.UUID, 0, len(members))
		for _, member := range members {
			userIDs = append(userIDs, member.UserID)
		}
		return userIDs, nil
	default:
		return nil, notificationsDomainNotifications.ErrInvalidAudienceType
	}
}
