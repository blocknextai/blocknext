# Web3

> Ethereum signature verification used for Metamask login. Shared primitive, not a full bounded context.

## Responsibility
This module owns Ethereum signature verification for the Metamask login flow: building the login message and verifying/recovering the signer of a personal-signed message. It has no persistence, no HTTP routes, and no aggregates — it only exposes the `SignatureVerifier` consumed by the `account` context.

## Domain
- **Key types:**
  - `crypto.SignatureVerifier` — `VerifySignature` / `PersonalEcRecover` / `BuildLoginMessage`.

## Dependencies
- **Bounded contexts:** none; provides `SignatureVerifier` to `account` (Metamask login).
- **Infrastructure:** none — no DB, cache, or outbound clients. The login message template comes from `AuthMetamaskOptions.LoginMessage` via `Dependencies.LoginMessage`.

## Layout
- `domain/crypto/` — `SignatureVerifier` interface + errors.
- `infrastructure/crypto/` — Ethereum implementation (`NewEthereumSignatureVerifier`).
- `module.go` — wires the verifier and exposes it on `Module`.
