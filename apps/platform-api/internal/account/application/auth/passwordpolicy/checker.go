package passwordpolicy

import (
	"context"
	"crypto/sha1"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/hex"
	"github.com/blocknextai/go-packages/httpclient"
)

const (
	hibpBaseURL = "https://api.pwnedpasswords.com"
	hibpTimeout = 1 * time.Second
)

var (
	errHIBPUnavailable = errors.New("hibp service unavailable")
)

type checker struct {
	hibpClient *httpclient.Client
}

func NewChecker() Policy {
	return &checker{
		hibpClient: httpclient.NewClientBuilder().
			BaseURL(hibpBaseURL).
			Timeout(hibpTimeout).
			Header("Add-Padding", "true").
			Build(),
	}
}

func (c *checker) Check(ctx context.Context, password string, userInputs []string) error {
	breached, err := c.checkHIBP(ctx, password)
	if err != nil {
		slog.WarnContext(ctx, "HIBP check failed; failing open",
			"component", "passwordpolicy",
			"error", err)
		return nil
	}

	if breached {
		return ErrPasswordBreached
	}

	return nil
}

func (c *checker) checkHIBP(ctx context.Context, password string) (bool, error) {
	sum := sha1.Sum([]byte(password))
	hash := strings.ToUpper(hex.Encode(sum[:]))
	prefix := hash[:5]
	suffix := hash[5:]

	var path strings.Builder
	path.WriteString("/range/")
	path.WriteString(prefix)

	response, err := c.hibpClient.Get(path.String()).
		Context(ctx).
		DoRaw()
	if err != nil {
		return false, err
	}

	if !response.IsSuccess() {
		return false, errHIBPUnavailable
	}

	for _, rawLine := range strings.Split(string(response.Body), "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		if strings.EqualFold(parts[0], suffix) {
			return true, nil
		}
	}

	return false, nil
}
