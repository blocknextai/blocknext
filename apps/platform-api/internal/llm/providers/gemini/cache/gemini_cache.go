package cache

import (
	"context"
	"crypto/sha256"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/httpclient"
	llmCache "github.com/blocknextai/platform-api/internal/llm/cache"
	"github.com/blocknextai/platform-api/internal/llm/providers/gemini"
)

// TODO: make configurable
const (
	ttl              = time.Hour
	refreshThreshold = 5 * time.Minute
	deleteTimeout    = 10 * time.Second
)

type geminiCache struct {
	mu              sync.Mutex
	model           string
	component       string
	name            string
	instructionHash [32]byte
	expiresAt       time.Time
	httpClient      *httpclient.Client
}

func New(
	apiKey string,
	model string,
	component string,
	timeout time.Duration,
) llmCache.Cache {
	builder := httpclient.NewClientBuilder().
		BaseURL(gemini.BaseURL).
		QueryParam("key", apiKey)
	if timeout > 0 {
		builder.Timeout(timeout)
	}
	return &geminiCache{
		model:      model,
		component:  component,
		httpClient: builder.Build(),
	}
}

func (c *geminiCache) Ensure(ctx context.Context, systemInstruction string) string {
	hash := sha256.Sum256([]byte(systemInstruction))

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.name != "" && c.instructionHash == hash && time.Until(c.expiresAt) > refreshThreshold {
		return c.name
	}

	name, expiresAt, err := c.create(ctx, systemInstruction)
	if err != nil {
		slog.Warn("Gemini cached content creation failed; falling back to inline systemInstruction",
			"component", c.component,
			"model", c.model,
			"error", err,
		)
		if c.name != "" && c.instructionHash == hash {
			return c.name
		}
		c.name = ""
		c.instructionHash = [32]byte{}
		c.expiresAt = time.Time{}
		return ""
	}

	oldName := c.name
	c.name = name
	c.instructionHash = hash
	c.expiresAt = expiresAt

	slog.Info("Gemini cached content created",
		"component", c.component,
		"model", c.model,
		"name", name,
		"expires_at", expiresAt,
	)

	if oldName != "" {
		go c.deleteAsync(oldName)
	}

	return name
}

func (c *geminiCache) create(ctx context.Context, systemInstruction string) (string, time.Time, error) {
	var success cachedContentResponse
	var errResp gemini.ErrorResponse
	response, err := c.httpClient.Post("/cachedContents").
		Context(ctx).
		JSONContentType().
		Body(map[string]any{
			"model": "models/" + c.model,
			"system_instruction": map[string]any{
				"parts": []map[string]any{
					{
						"text": systemInstruction,
					},
				},
			},
			"ttl": strconv.FormatInt(int64(ttl.Seconds()), 10) + "s",
		}).
		Do(&success, &errResp)

	if err != nil {
		return "", time.Time{}, errCreateFailed.WithCause(err)
	}

	if !response.IsSuccess() {
		return "", time.Time{}, errCreateFailed.WithCause(apperror.Internal(errResp.Error.Message))
	}

	expiresAt := time.Now().Add(ttl)
	if parsed, parseErr := time.Parse(time.RFC3339, success.ExpireTime); parseErr == nil {
		expiresAt = parsed
	}

	return success.Name, expiresAt, nil
}

func (c *geminiCache) deleteAsync(name string) {
	ctx, cancel := context.WithTimeout(context.Background(), deleteTimeout)
	defer cancel()

	_, err := c.httpClient.Delete("/"+name).
		Context(ctx).
		Do(nil, nil)
	if err != nil {
		slog.Warn("Gemini cached content delete failed",
			"component", c.component,
			"name", name,
			"error", err,
		)
	}
}
