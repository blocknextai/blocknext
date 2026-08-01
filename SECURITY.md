# Security Policy

## Supported Versions

Security fixes are applied to the latest version on the `main` branch.

## Reporting a Vulnerability

Please **do not** open a public issue for security vulnerabilities.

Use [GitHub private vulnerability reporting](https://docs.github.com/en/code-security/security-advisories/guidance-on-reporting-and-writing-information-about-vulnerabilities/privately-reporting-a-security-vulnerability) instead:

1. Go to the repository's **Security** tab.
2. Click **Report a vulnerability**.
3. Describe the issue, the affected component (`apps/platform-api`, `apps/file-gateway-api`, `apps/platform`, `packages/go-packages`), and reproduction steps if possible.

We will respond through the advisory thread. Please give us a reasonable amount of time to address the issue before any public disclosure.

## Scope

- Authentication/authorization bypasses
- Injection of any kind (SQL, command, template, ...)
- Sensitive data exposure
- SSRF, XSS, CSRF and similar web vulnerabilities
- Vulnerabilities in the file upload/download pipeline

Self-hosted deployment misconfigurations (e.g. exposing the database port publicly, running with example secrets) are out of scope.
