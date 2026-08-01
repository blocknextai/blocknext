package publishpost

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/x/helpers"
)

type XPublishPostExecutorInput struct {
	Text           string `schema:"text"`
	ReplyToTweetID string `schema:"replyToTweetId"`
}

type XPublishPostExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[XPublishPostExecutorInput]
}

func NewXPublishPostExecutor(
	nodeID string,
	validator *jsonschema.Validator[XPublishPostExecutorInput],
) *XPublishPostExecutor {
	return &XPublishPostExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

func (e *XPublishPostExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "x_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			client := helpers.GetXClient(ctx, accessToken)

			tweetData := map[string]any{
				"text": input.Text,
			}

			if input.ReplyToTweetID != "" {
				tweetData["reply"] = map[string]any{
					"in_reply_to_tweet_id": input.ReplyToTweetID,
				}
			}

			var errorResponse helpers.ErrorResponse

			response, err := client.Post("/2/tweets").
				Body(tweetData).
				Do(nil, &errorResponse)

			rateLimitInfo := helpers.ExtractRateLimitInfo(response.Headers)

			result := map[string]any{
				"rateLimit": helpers.RateLimitInfoToMap(rateLimitInfo),
			}

			if err != nil {
				result["status"] = false
				results = append(results, result)
				return results, err
			}

			if !response.IsSuccess() {
				result["status"] = false
				results = append(results, result)

				return results, apperror.Internal(errorResponse.Detail)
			}

			result["status"] = true
			results = append(results, result)
		}

		return results, nil
	}
}
