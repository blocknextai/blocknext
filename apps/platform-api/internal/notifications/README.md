# Notifications

> Event-driven inbox that fans out domain events from other contexts into per-user, in-app notifications.

## Responsibility
This context owns the in-app notification inbox. It subscribes to domain events published across the system, persists a broadcast `Notification`, and fans it out into per-user `NotificationRecipient` rows resolved by audience (a single user, or every member of an organization). It also serves the read/manage side of the inbox (list, counts, mark read/seen, delete). It does not send email/push and publishes no events of its own.

## Domain
- **Aggregates / key types:**
  - `notifications.Notification` — the broadcast record (`Type`, `Level`, `AudienceType`, `AudienceID`, `Title`, `Body`, `Data`, `ActionURL`).
  - `notificationrecipients.NotificationRecipient` — one user's copy/state of a notification (`ReadAt`, `SeenAt`; soft-delete).
  - `notificationrecipients.InboxItem` — read-model join of recipient + notification for inbox listing.
  - `Level` (`info`/`success`/`warning`/`error`) and `AudienceType` (`user`/`organization`) enums.
- **Key rules / invariants:**
  - `Notification` requires non-blank `Type`/`Title`, valid `Level`/`AudienceType`, and non-nil `AudienceID`.
  - `NotificationRecipient.MarkRead()` is idempotent and also sets `SeenAt` if unset.

## Use cases (application)
- **`application/notifications`** — `NotificationService.Create` persists the notification then resolves recipients by audience (`organization` audience fans out via `OrganizationUserService`) and bulk-inserts recipient rows.
- **Queries:** `getallnotifications` (paginated inbox), `getnotificationcounts` (unread/unseen).
- **Commands:** `marknotificationread`, `markallnotificationsread`, `markallnotificationsseen`, `deletenotification`.

## HTTP API
All routes require `Authenticate()` plus an RBAC permission. Mirrored under a user scope and an organization scope.

User base `/users/me/notifications` (`RequireUserPermission`):
- `GET /` — list inbox (`ReadNotificationPermission`)
- `GET /count` — unread/unseen counts (`ReadNotificationPermission`)
- `POST /seen` — mark all seen (`UpdateNotificationPermission`)
- `POST /read-all` — mark all read (`UpdateNotificationPermission`)
- `PATCH /:recipientId/read` — mark one read (`UpdateNotificationPermission`)
- `DELETE /:recipientId` — delete one (`DeleteNotificationPermission`)

Organization base `/organizations/:organizationId/notifications` — same methods/paths/permissions via `RequireOrganizationPermission`.

## Events
- **Published:** none.
- **Consumed:** 5 eventbus subscribers (each idempotent via the inbox dedupe key `notifications:<event_name>`, calling `NotificationService.Create`):
  - `account.user.created` → welcome notification.
  - `organizations.organization_user.created` → "added to organization" (skipped for the owner role).
  - `organizations.organization_user.role_changed` → role-changed notice.
  - `account.user.password_changed` / `account.user.email_changed` → security notices.

## Dependencies
- **Bounded contexts:** organizations (`OrganizationUserService`, to expand organization audiences).
- **Infrastructure:** Postgres schema `notifications` (tables `notifications.notifications`, `notifications.recipients`); eventbus + inbox idempotency service. No external APIs.

## Layout
Standard DDD layers present, wired in `module.go`.
