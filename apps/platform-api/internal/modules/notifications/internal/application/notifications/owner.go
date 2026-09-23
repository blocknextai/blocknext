package notifications

import (
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	notificationsDomainNotifications "github.com/blocknextai/platform-api/internal/modules/notifications/internal/domain/notifications"
)

func AudienceTypeFromOwnerType(ownerType commonDomain.OwnerType) notificationsDomainNotifications.AudienceType {
	if ownerType == commonDomain.OwnerTypeOrganization {
		return notificationsDomainNotifications.AudienceTypeOrganization
	}
	return notificationsDomainNotifications.AudienceTypeUser
}
