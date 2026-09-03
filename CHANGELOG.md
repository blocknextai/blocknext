# Changelog

## [0.0.10](https://github.com/blocknextai/blocknext/releases/tag/v0.0.10) - 2026-09-03

### Features
- **tour:** Cover the AI chat panel and follow the user's clicks ([dc081a6](https://github.com/blocknextai/blocknext/commit/dc081a676f69e57166b35509bf68edb7bf4d18e8))
- **tour:** Walk the onboarding tour through a real flow ([fef0488](https://github.com/blocknextai/blocknext/commit/fef048884c4dc8bb42e16e0bf3e2cc708d1887a1))

### Refactoring
- **config:** Embed the system instruction prompts ([067fad2](https://github.com/blocknextai/blocknext/commit/067fad205fd706de7aeddaa189ac70eccec5d6f5))

### Documentation
- Add badges to the README ([f7a3e63](https://github.com/blocknextai/blocknext/commit/f7a3e63d8e3d434e53243206c51057f22e16a765))

### Maintenance
- Generate CHANGELOG.md on release with git-cliff ([b86ed3c](https://github.com/blocknextai/blocknext/commit/b86ed3c15151dd26c89b87414b993023c58176a9))

## [0.0.9](https://github.com/blocknextai/blocknext/releases/tag/v0.0.9) - 2026-08-25

### Features
- **mcp:** Gate, record and harden tool calls ([521363f](https://github.com/blocknextai/blocknext/commit/521363fd5649425b4cd76a537f3703c01f6dfe90))
- **mcp:** Gate, record and harden tool calls ([eaef7f3](https://github.com/blocknextai/blocknext/commit/eaef7f3bc3e64126d213dde57e074b079eadd4a5))

### Bug Fixes
- **platform-api:** Wire CredentialService into the credential OAuth module ([46a8af2](https://github.com/blocknextai/blocknext/commit/46a8af28d1afa1209338e64c99eaf9530a8c65e1))

### Refactoring
- **BREAKING:** Remove MetaMask login ([939fcbf](https://github.com/blocknextai/blocknext/commit/939fcbfa1149f82f1494fdd48e692809360b2175))
- **platform:** Centralize login redirect and logout handling ([1554892](https://github.com/blocknextai/blocknext/commit/155489254bea4d42362b31855380beb1451d136d))

### Maintenance
- Lint on bun 1.4.0 ([3bb5613](https://github.com/blocknextai/blocknext/commit/3bb5613a18d473f296f4b65d7dae5495c73d77a2))
- Enable weekly dependabot updates ([4be296b](https://github.com/blocknextai/blocknext/commit/4be296b5654451fdad36dcff395812a4906d3110))
- Rename the local docker workflow to docker-dev ([cfcb50a](https://github.com/blocknextai/blocknext/commit/cfcb50a7dd37a8718b3112abd9df9950b00e2779))

## [0.0.8](https://github.com/blocknextai/blocknext/releases/tag/v0.0.8) - 2026-08-20

### Features
- Add a PowerShell setup script and per-OS install steps ([4d15b1c](https://github.com/blocknextai/blocknext/commit/4d15b1c4c23f131f8169b3d77a89913c8827308d))
- **flow-editor:** Offer every upstream node as a data source ([da9a3f9](https://github.com/blocknextai/blocknext/commit/da9a3f911c0ac3accbf02e682155eefb4620389b))

### Documentation
- Add flow editor screenshot to README ([e43c19c](https://github.com/blocknextai/blocknext/commit/e43c19c399a915abf91a1d765f84773c54aeccc7))

### Maintenance
- Upgrade to Go 1.27 and bump Go dependencies ([8bac39a](https://github.com/blocknextai/blocknext/commit/8bac39aff83178a3a31d74986653e528468394e0))
- Upgrade to Go 1.27 and bump Go dependencies ([1dd8a97](https://github.com/blocknextai/blocknext/commit/1dd8a97c759a3d154ad6f9b28b057af9e6c1a9c2))

## [0.0.7](https://github.com/blocknextai/blocknext/releases/tag/v0.0.7) - 2026-08-10

### Features
- **taskrunner:** Route items per branch and harden the run ([7725203](https://github.com/blocknextai/blocknext/commit/7725203bb12c01e63d2ec35b1526ba549c03d34a))

### Bug Fixes
- **flow:** Keep a renamed node and widen what a node can read ([1f1c14e](https://github.com/blocknextai/blocknext/commit/1f1c14e70145caf57077c7f23043147738ddff7a))

### Refactoring
- **ui:** Render images through one component ([7680374](https://github.com/blocknextai/blocknext/commit/768037429e37e33332c04cc3283112d91c41f9de))
- **taskrunner:** Give the reference language one substitution path ([d2c04b1](https://github.com/blocknextai/blocknext/commit/d2c04b1f9a566e6a10eedcb2d6190fe28903bfae))

## [0.0.6](https://github.com/blocknextai/blocknext/releases/tag/v0.0.6) - 2026-08-09

### Features
- **mcp:** Show how to connect to a server ([ecd5476](https://github.com/blocknextai/blocknext/commit/ecd5476a6a4a1630cf1990bbbb4f574dd09a0b92))
- **flow:** Move annotation into the node catalog and fix canvas view ([d68de78](https://github.com/blocknextai/blocknext/commit/d68de78c9ac682291c51c4916e4b65e856585d0d))

### Bug Fixes
- **flow:** Make connection points and their labels readable ([86e7783](https://github.com/blocknextai/blocknext/commit/86e7783f3fdda10e1da4dcaeb1c6d888c6d260c1))
- **flow:** Stop defaulting nodes to a 30s timeout ([dac18b6](https://github.com/blocknextai/blocknext/commit/dac18b66b4ee2a874b5b6af768ed53e048100fcd))
- **compose:** Make the default self-hosted install durable ([c83f2a0](https://github.com/blocknextai/blocknext/commit/c83f2a0d09576a86c8a3246b8d4e5854cf044b3f))

### Documentation
- Describe the batch model instead of "no loops" ([99955a7](https://github.com/blocknextai/blocknext/commit/99955a734afde9f3beb8fbd40f2575ee6540d50c))

## [0.0.5](https://github.com/blocknextai/blocknext/releases/tag/v0.0.5) - 2026-08-07

### Features
- **flow-editor:** Render nodes from what the API declares ([565dac9](https://github.com/blocknextai/blocknext/commit/565dac9184ff4d4d7a5a14405fda6d606ae7337e))
- **platform:** Ship brand marks and action glyphs as files ([999d5dd](https://github.com/blocknextai/blocknext/commit/999d5dd28f507e1861ace31fdc9cb8593c73b7d7))
- **nodeengine:** Let a node state its own icon and connection points ([8a7aa29](https://github.com/blocknextai/blocknext/commit/8a7aa295c89a87c6e3d6b31b8683ed4f52a783a1))
- **infra:** Make Redis optional and default to in-process providers ([c284595](https://github.com/blocknextai/blocknext/commit/c284595131e10f4f96946ad2c42de3e5e6263b2b))
- **platform:** Expose the AI feature flags and gate the UI surfaces ([d130b27](https://github.com/blocknextai/blocknext/commit/d130b27ea24f471566cc6d46eb03ce2a35cfbb92))

### Bug Fixes
- **flow-editor:** Skip the save request when running an unchanged flow ([bd88659](https://github.com/blocknextai/blocknext/commit/bd886590abf617d6262d3465841c10e67e076589))

### Refactoring
- **dag:** Record which output an edge leaves from ([27fd6e4](https://github.com/blocknextai/blocknext/commit/27fd6e43df224a60968a09848a356cd548ef5711))

### Documentation
- **nodeengine:** Describe how a node states its icon and handles ([daf329d](https://github.com/blocknextai/blocknext/commit/daf329dffc23a137abfd94234630e16770cf07f7))

### Maintenance
- **flow-editor:** Widen the API trigger sheet ([368ad76](https://github.com/blocknextai/blocknext/commit/368ad767355c52d327597e916b9dff21d38cfebf))
- **env:** Enable the log email sender by default ([17b450b](https://github.com/blocknextai/blocknext/commit/17b450b8b80778483cb37206b1702fa7eb2360ed))
- **docker:** Start queue-only services only in queue mode ([1c59eb7](https://github.com/blocknextai/blocknext/commit/1c59eb7a24c3323bdf456b9b17cf4c1f7ba295ae))
- **go:** Update dependencies and pin the in-repo go-packages requirement ([894ff23](https://github.com/blocknextai/blocknext/commit/894ff2303103f8e36dd21b73673a0e2e4c0a070a))

## [0.0.4](https://github.com/blocknextai/blocknext/releases/tag/v0.0.4) - 2026-08-03

### Features
- **triggers:** Manage the webhook secret and label webhook sources ([5480032](https://github.com/blocknextai/blocknext/commit/5480032227310ae5432ccb48fa5267da3872aa69))
- **llm:** Log Gemini token usage and drop function-calling context cache ([ec25c74](https://github.com/blocknextai/blocknext/commit/ec25c7413c15d2f34339a107aa6aa5b4f54cb6a5))

### Bug Fixes
- **credentials:** Close popover and reset query on secret selection ([c7f40d4](https://github.com/blocknextai/blocknext/commit/c7f40d4cb9809c75845d4c07484eb28872d2b5c3))

### Documentation
- Add async-encapsulation and featherweight self-hosting highlights ([24ac065](https://github.com/blocknextai/blocknext/commit/24ac06520fe486f0db3909175c5dfa3c2224ab8a))

### Maintenance
- Move community health files under .github/ ([ac602b6](https://github.com/blocknextai/blocknext/commit/ac602b6d8a99926d4fa0a79c5d3bb678a5ba9894))

## [0.0.3](https://github.com/blocknextai/blocknext/releases/tag/v0.0.3) - 2026-08-03

### Maintenance
- Split setup into env-only setup and setup-dev ([8a8946f](https://github.com/blocknextai/blocknext/commit/8a8946f3ad8afd5899620457e5f530d4d1eab3b4))

## [0.0.2](https://github.com/blocknextai/blocknext/releases/tag/v0.0.2) - 2026-08-02

### Documentation
- Enrich module docs, add architecture overview and community templates ([552e00c](https://github.com/blocknextai/blocknext/commit/552e00c016faa84fa79fd9b37904210f4290be30))
- Update stale READMEs after OSS cleanup ([dd78146](https://github.com/blocknextai/blocknext/commit/dd78146593f067a738961a3eddac8ec4633c1779))

## [0.0.1](https://github.com/blocknextai/blocknext/releases/tag/v0.0.1) - 2026-08-01

### Features
- Initial open-source release ([595ad84](https://github.com/blocknextai/blocknext/commit/595ad846cbde99974e3d6638bf224d5fc8eb6977))

