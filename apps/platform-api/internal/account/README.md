# Account

> Owns user identity, authentication, and per-user account settings across password, email/magic-link, OAuth, and crypto-wallet sign-in.

## Responsibility
This context owns the `User` aggregate and everything that authenticates a person and links them to login identities: password credentials, OAuth/crypto linked accounts, sessions with refresh-token rotation, OAuth nonces/PKCE, and one-time verification tokens. It also owns user-facing account settings (preferences, social links) and emits domain events that drive its own transactional emails. It does NOT own organizations or RBAC role/permission definitions (it consumes the common `UserPermissionChecker`).

## Domain
- **Aggregates / key types:**
  - `User` — the account identity, holding role, verification, and ban status.
  - `Session` — an authenticated session with refresh-token rotation (`TokenFamily`/`TokenGeneration`) and revocation.
  - `LinkedAccount` — a mapping from a user to an external provider identity (email, OAuth, crypto); one primary per user.
  - `PasswordCredential` — bcrypt password hash for email/password auth (one per user).
  - `UserNonce` — OAuth2 nonce + PKCE verifier/challenge with short TTL.
  - `VerificationToken` — one-time, hash-only token with a purpose (`email_verify`, `password_reset`, `email_change`, `magic_link`).
  - `UserPreference` — theme/language settings.
  - `UserSocial` — ordered social-media links.
  - Provider taxonomy: `AuthProviderType` (`metamask`, `google`, `x`, `facebook`, `github`, `email`, `password`) and `AuthProviderCategory` (`crypto`, `oauth`, `email`, `password`).
- **Key rules / invariants:**
  - Tokens and password hashes are stored hashed only; plaintext is never persisted.
  - Verification tokens and nonces are expiry-checked and soft-deleted on consume.
  - One primary linked account per user; `ProviderID` unique; emails normalized.
  - Refresh tokens rotate by family/generation to detect reuse.

## Use cases (application)
**Commands**
- Auth (OAuth/crypto): `createusernonce` (start OAuth + PKCE), `createusertoken` (exchange code for session), `refreshtoken` (rotate access token).
- Auth (password): `password/register`, `password/login`, `password/forgot`, `password/reset`, `password/change`, `password/set` (initial password on OAuth user), `password/verify` (email verify).
- Auth (email): `email/add`, `email/change`, `email/confirm`, `email/resendverification`.
- Auth (magic link): `email/magiclink/request`, `email/magiclink/consume`.
- Sessions: `logout`, `revokesession`, `revokeallsessions`.
- Linked accounts: `addlinkedaccount`, `deletelinkedaccount`.
- User settings: `userpreferences/updateuserpreferences`, `usersocials/updateusersocial`.

**Queries**
- `auth/getauthmethods` — available auth providers/settings.
- `users/getprofile`, `users/getroles`.
- `sessions/getallsessions`, `linkedaccounts/getalllinkedaccounts`.
- `userpreferences/getuserpreferences`, `usersocials/getallusersocials`.

## HTTP API
**Auth (public)**
| Method | Path | Purpose |
| - | - | - |
| GET | `/auth/methods` | List auth methods (cached) |
| POST | `/auth/oauth/nonce` | Start OAuth flow (nonce + PKCE) |
| POST | `/auth/oauth/token` | Exchange OAuth code for tokens |
| POST | `/auth/token/refresh` | Refresh access token |
| POST | `/auth/password/register` | Register with email/password |
| POST | `/auth/password/login` | Password login |
| POST | `/auth/password/forgot` | Request password reset |
| POST | `/auth/password/reset` | Complete password reset |
| POST | `/auth/email/verify` | Verify email |
| POST | `/auth/email/resend-verification` | Resend verification |
| POST | `/auth/email/change/confirm` | Confirm email change |
| POST | `/auth/email/magic-link/request` | Request magic link |
| POST | `/auth/email/magic-link/consume` | Consume magic link |

**Authenticated**
| Method | Path | Purpose | Auth note |
| --- | --- | --- | --- |
| POST | `/auth/logout` | Revoke current session | Authenticate |
| GET | `/users/me/` | Current user profile | `ReadUserPermission` |
| GET | `/users/me/roles` | User roles | `ReadUserPermission` |
| POST | `/users/me/email/add` | Add email | Authenticate |
| POST | `/users/me/email/change` | Request email change | Authenticate |
| POST | `/users/me/password/set` | Set initial password | Authenticate |
| POST | `/users/me/password/change` | Change password | Authenticate |
| GET / POST | `/users/me/linked-accounts/` | List / link account | `Read`/`CreateLinkedAccount(s)Permission` |
| DELETE | `/users/me/linked-accounts/:linkedAccountId` | Unlink account | `DeleteLinkedAccountPermission` |
| GET / PUT | `/users/me/socials/` | List / update socials | `Read`/`UpdateSocialPermission` |
| GET / PATCH | `/users/me/preferences/` | Get / update preferences | `Read`/`UpdateUserPreferencesPermission` |
| GET | `/users/me/sessions/` | List sessions | `ReadSessionPermission` |
| POST | `/users/me/sessions/revoke-all` | Revoke all sessions | `DeleteSessionPermission` |
| POST | `/users/me/sessions/:sessionId/revoke` | Revoke one session | `DeleteSessionPermission` |

## Events
- **Published:** `account.user.created`, `account.user.email_changed`, `account.user.password_changed`, `account.email.added`, `account.email_verification.requested`, `account.email_change.requested`, `account.password_reset.requested`, `account.magic_link.created`, `account.registration.existing_email_notified`.
- **Consumed (self):** the email-driving events above plus `account.user.created` are consumed by handlers under `application/events/` to send the corresponding transactional emails (welcome, verification, change-confirm, reset, magic link, existing-email notice).

## Dependencies
- **Bounded contexts:** None consumed via service interfaces. Exposes `SessionService`, `UserService`, `LinkedAccountService`, and the `UserPermissionChecker` for other modules/middleware.
- **Infrastructure:** Postgres tables in schema `account` — `users`, `linked_accounts`, `nonces`, `socials`, `user_preferences`, `sessions`, `password_credentials`, `verification_tokens`. OAuth providers (Google, GitHub, X, Facebook) and Metamask crypto-signature verification (`web3` `SignatureVerifier`). JWT service, email sender, password hasher, cache, and the eventbus (publish + inbox idempotency).

## Layout
Standard DDD layers: `domain/` (aggregates + rules), `application/` (CQRS handlers + event handlers), `infrastructure/` (persistence + auth-provider adapters), `presentation/` (HTTP). Wired in `module.go`.
