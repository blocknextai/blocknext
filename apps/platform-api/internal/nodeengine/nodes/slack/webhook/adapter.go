package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"strconv"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/hex"
	"github.com/blocknextai/go-packages/json"
	nodeEngineDomainAdapters "github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
)

const (
	slackSignatureMaxAge = 5 * time.Minute
)

type SlackAdapter struct {
	nodeEngineDomainAdapters.Adapter
}

func NewSlackAdapter(nodeID string) *SlackAdapter {
	return &SlackAdapter{
		Adapter: nodeEngineDomainAdapters.Adapter{
			ID:   nodeID,
			Name: "Slack",
		},
	}
}

type slackPayload struct {
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
	Event     struct {
		Text string `json:"text"`
		User string `json:"user"`
	} `json:"event"`
}

func (a *SlackAdapter) Adapt(raw map[string]any) (*nodeEngineDomainAdapters.TriggerContext, error) {
	var payload slackPayload
	if err := json.ArgsToStruct(raw, &payload); err != nil {
		return nil, err
	}

	return &nodeEngineDomainAdapters.TriggerContext{
		Source:  a.ID,
		Sender:  payload.Event.User,
		Prompt:  payload.Event.Text,
		Payload: raw,
	}, nil
}

var (
	ErrSlackVerificationFailed = apperror.Internal("verification failed")
	ErrSlackMissingSignature   = apperror.Unauthorized("missing slack signature headers")
	ErrSlackStaleSignature     = apperror.Unauthorized("slack signature timestamp too old")
	ErrSlackInvalidSignature   = apperror.Unauthorized("invalid slack signature")
)

func (a *SlackAdapter) Verify(request *nodeEngineDomainAdapters.VerificationRequest, _ *string) (*nodeEngineDomainAdapters.VerificationResponse, error) {
	var payload slackPayload
	if err := json.Unmarshal(request.Body, &payload); err != nil {
		return nil, err
	}

	if payload.Type != "url_verification" {
		return nil, nil
	}

	if payload.Challenge == "" {
		return nil, ErrSlackVerificationFailed
	}

	return &nodeEngineDomainAdapters.VerificationResponse{
		Body:       []byte(payload.Challenge),
		StatusCode: 200,
	}, nil
}

func (a *SlackAdapter) ValidateSignature(request *nodeEngineDomainAdapters.VerificationRequest, secret *string) error {
	if secret == nil || strings.TrimSpace(*secret) == "" {
		return nil
	}

	signature := request.Header("X-Slack-Signature")
	timestamp := request.Header("X-Slack-Request-Timestamp")
	if strings.TrimSpace(signature) == "" || strings.TrimSpace(timestamp) == "" {
		return ErrSlackMissingSignature
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return ErrSlackInvalidSignature
	}
	if time.Since(time.Unix(ts, 0)).Abs() > slackSignatureMaxAge {
		return ErrSlackStaleSignature
	}

	mac := hmac.New(sha256.New, []byte(*secret))
	mac.Write([]byte("v0:" + timestamp + ":"))
	mac.Write(request.Body)
	expected := "v0=" + hex.Encode(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return ErrSlackInvalidSignature
	}

	return nil
}
