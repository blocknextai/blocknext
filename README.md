# BlockNext

BlockNext is a no-code platform for building and running AI-powered workflows. Design flows on a drag-and-drop canvas, connect AI models and third-party services, and let your automations run — no code required.

It is built for end users, not just developers: flows stay simple, readable and predictable by design.

## Highlights

- **Visual flow builder** — compose workflows on an intuitive drag-and-drop canvas. Flows are deliberately simple: no loops, no sub-workflows, no expression language to learn.
- **Describe it, don't configure it** — give a node plain-language instructions and an LLM fills in its parameters at run time. Anything you set explicitly always wins over what the model infers.
- **Built-in MCP server** — every integration node doubles as an [MCP](https://modelcontextprotocol.io) tool. Point Claude (or any MCP client) at your BlockNext server with an API key and use your connected services from chat.
- **AI-powered nodes** — LLMs and generative AI (text, image, audio, video) as first-class building blocks, alongside integrations for the tools you already use.
- **Triggers** — start flows manually, on a schedule, or from the outside via webhooks and API calls.
- **Live execution view** — watch every task and node progress in real time over WebSocket.
- **Credentials management** — encrypted at rest, OAuth tokens auto-refreshed, and only ever decrypted at execution time — flows never embed secrets.

Curious how it works under the hood? See the [architecture overview](ARCHITECTURE.md).

## Getting started

```bash
make setup      # creates .env with generated secrets
make docker-up  # pulls the published images and starts the full stack
```

Docker and `make` are all you need — Go and Bun are only required for
[development](.github/CONTRIBUTING.md).

The UI is served on http://localhost:4000. Run `make help` for every target.

To build and run everything from source instead, see
[CONTRIBUTING.md](.github/CONTRIBUTING.md) — the development workflow uses
`make local-docker-up`.

## Community

- [Architecture overview](ARCHITECTURE.md)
- [Contributing guide](.github/CONTRIBUTING.md)
- [Code of Conduct](.github/CODE_OF_CONDUCT.md)
- [Security policy](.github/SECURITY.md) — please report vulnerabilities privately

## License

BlockNext is licensed under the [Apache License 2.0](LICENSE).
