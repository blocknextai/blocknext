---
name: upload-create-rule
description: Adds a new UploadRule to internal/upload/domain/upload/rules.go in the file-gateway-api project. Use when the user asks to "add a new upload type", "create an upload rule", "register a new upload", or otherwise wants to define a new entry in the upload rules registry.
---

> **Working directory:** all relative paths in this skill are relative to `apps/file-gateway-api/` in the monorepo.


# upload-create-rule

Add a new `UploadRule` to `internal/upload/domain/upload/rules.go` and register it in the `uploadRules` map.

## When to use

- User wants to register a new file upload type (new feature, new node, new asset category).
- User asks to extend marketplace, organization, api.nodes, api.web3, bug_reports, etc. with a new upload target.

## Data model

`internal/upload/rule.go` defines `UploadRule`:

| Field | Type | Notes                                         |
|------|-----|-----------------------------------------------|
| `ID` | `string` | UUID v7, must be unique across rules          |
| `Title` | `string` | Dot notation (e.g. `marketplace.cover_image`) |
| `MaxSize` | `int64` | Bytes — written as `N * 1024 * 1024`          |
| `AllowedMimes` | `[]string` | One `.AllowedMime(...)` call per entry        |
| `DefaultFolder` | `string` | Storage folder, must end with `/`             |
| `IsOverrideFilename` | `bool` | Set via `.OverrideFilename()`                 |
| `IsPublic` | `bool` | Set via `.Public()` or `.Private()`           |

## Steps

### 1. Gather inputs

If any of the following are missing from the user's request, ask via `AskUserQuestion` (one question per missing field, batched in a single call):

1. **Title** — e.g. `marketplace.video_thumbnail`, `api.nodes.openai.tts`.
2. **Max size in MB** — converted to `N * 1024 * 1024` in code.
3. **Allowed MIME category** — image, video, audio, document, json, or custom list.
4. **Default folder** — e.g. `marketplace/video-thumbnails/` (must end with `/`, kebab-case).
5. **Visibility** — public or private.
6. **Override filename** — usually yes for public uploads, no for private payloads.

### 2. Generate UUID

```bash
uuidgen | tr '[:upper:]' '[:lower:]'
```

Verify the UUID is not already present:

```bash
grep -i "<new-uuid>" /Users/bilal/workspaces/blocknext/projects/file-gateway-api/internal/upload/domain/upload/rules.go
```

### 3. Derive variable name

Convert the title to camelCase, following existing conventions in `rules.go`:

- `marketplace.cover_image` → `marketplaceCoverImage`
- `api.nodes.gemini.imagen` → `apiImagenNode`
- `api.web3.tokens.icons` → `apiWeb3TokenIcon`
- `api.organization.icon` → `organizationIconImage`
- `bug_reports.files` → `bugReportFiles`

Check for variable name collisions in `rules.go` before proceeding.

### 4. MIME presets

Reuse these sets from existing rules unless the user specifies otherwise:

**Image (full):**
```
image/png, image/jpeg, image/jpg, image/avif, image/heic, image/heif,
image/webp, image/gif, image/heif-sequence, image/heic-sequence
```

**Image (minimal):** `image/png`, `image/jpeg`, `image/jpg`

**Video:** `video/mp4`, `video/webm`, `video/ogg`

**Audio:** `audio/mpeg`, `audio/mp3`, `audio/mp4`, `audio/m4a`, `audio/m4b`, `audio/m4p`, `audio/m4v`

**Document:** `application/pdf`, `text/plain`, `application/msword`, `application/vnd.openxmlformats-officedocument.wordprocessingml.document`, `application/vnd.ms-excel`

**JSON:** `application/json`

### 5. Insert the rule

Add the new `var` to `internal/upload/domain/upload/rules.go` near logically-related rules (e.g. a new `marketplace.*` rule goes next to other marketplace rules). Mirror the existing builder pattern exactly:

```go
var <variableName> = NewUploadRuleBuilder().
    ID("<uuid>").
    Title("<title>").
    MaxSize(<N> * 1024 * 1024).
    AllowedMime("<mime1>").
    AllowedMime("<mime2>").
    DefaultFolder("<folder>/").
    OverrideFilename().    // omit for private payloads
    Public().              // or Private()
    Build()
```

### 6. Register in the map

Append an entry to `uploadRules` at the bottom of the file:

```go
var uploadRules = map[string]UploadRule{
    // ...existing entries
    <variableName>.ID: <variableName>,
}
```

Don't worry about column alignment — `gofmt` handles it.

### 7. Verify

Always run a build after editing:

```bash
cd /Users/bilal/workspaces/blocknext/projects/file-gateway-api && go build ./...
```

If the build fails, fix the issue before reporting completion. Optionally run `go fmt ./internal/upload/...` to lock in formatting.

## Conventions (observed in the existing codebase)

- **Public rules** typically chain `.OverrideFilename()` — public assets need predictable, collision-free filenames.
- **Private rules** (e.g. webhook payloads) typically preserve the original filename — `.OverrideFilename()` is omitted.
- **Folder paths** end with `/` and use kebab-case for multi-word segments (`get-file/`, `image-generation/`).
- **Title namespaces** use dot notation: `<area>.<feature>.<subfeature>`. Existing areas: `marketplace`, `api.nodes.<provider>`, `api.web3`, `api.organization`, `bug_reports`.
- **Common MaxSize values:** `2 * 1024 * 1024` (images), `5 * 1024 * 1024` (documents/payloads), `10 * 1024 * 1024` (audio/video/large media).

## Output

Report concisely to the user when done:

- Generated UUID
- Variable name added
- File and approximate line range edited
- Build result (success / failure)
