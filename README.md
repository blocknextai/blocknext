# BlockNext

BlockNext is a no-code platform for building and running AI-powered workflows. Design flows on a drag-and-drop canvas, connect AI models and third-party services, and let your automations run — no code required.

It is built for end users, not just developers: flows stay simple, readable and predictable by design.

## Highlights

- **Visual flow builder** — compose workflows on an intuitive drag-and-drop canvas.
- **AI-powered nodes** — bring LLMs and generative AI into your flows as first-class building blocks.
- **Integrations** — connect the tools and services you already use.
- **Triggers** — start flows manually, on a schedule, or from the outside via webhooks and API calls.
- **Credentials management** — store and reuse service credentials securely across flows.

## Getting started

```bash
make setup      # creates .env with generated secrets, installs dependencies
make docker-up  # pulls the published images and starts the full stack
```

The UI is served on http://localhost:4000. Run `make help` for every target.

To build and run everything from source instead, see
[CONTRIBUTING.md](CONTRIBUTING.md) — the development workflow uses
`make local-docker-up`.

## Community

- [Contributing guide](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Security policy](SECURITY.md) — please report vulnerabilities privately

## License

BlockNext is licensed under the [Apache License 2.0](LICENSE).
