package publishpost

import (
	"context"
	"net/url"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/linkedin/helpers"
)

type LinkedinPublishPostExecutorInput struct {
	Text       string `schema:"text"`
	Visibility string `schema:"visibility"`
}

type LinkedinPublishPostExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[LinkedinPublishPostExecutorInput]
}

type LinkedinProfileResponse struct {
	ID string `json:"id"`
}

type LinkedinProfileError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type LinkedinPostRequest struct {
	Author          string `json:"author"`
	LifecycleState  string `json:"lifecycleState"`
	SpecificContent struct {
		ShareContent struct {
			ShareCommentary struct {
				Text string `json:"text"`
			} `json:"shareCommentary"`
			ShareMediaCategory string `json:"shareMediaCategory"`
		} `json:"com.linkedin.ugc.ShareContent"`
	} `json:"specificContent"`
	Visibility struct {
		MemberNetworkVisibility string `json:"com.linkedin.ugc.MemberNetworkVisibility"`
	} `json:"visibility"`
}

type LinkedinPostResponse struct {
	ID string `json:"id"`
}

type LinkedinPostError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewLinkedinPublishPostExecutor(
	nodeID string,
	validator *jsonschema.Validator[LinkedinPublishPostExecutorInput],
) *LinkedinPublishPostExecutor {
	return &LinkedinPublishPostExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

func (e *LinkedinPublishPostExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "linkedin_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			postID, err := e.createPost(client, *input)
			if err != nil {
				return nil, err
			}

			postURL := ""
			if postID != "" {
				postURL = "https://www.linkedin.com/feed/update/" + url.PathEscape(postID) + "/"
			}

			results = append(results, map[string]any{
				"status":  true,
				"postId":  postID,
				"postUrl": postURL,
			})
		}

		return results, nil
	}
}

func (e *LinkedinPublishPostExecutor) createPost(client *httpclient.Client, input LinkedinPublishPostExecutorInput) (string, error) {
	var profileResp LinkedinProfileResponse
	var profileErr LinkedinProfileError

	profileRes, err := client.Get("/people/~").
		JSONContentType().
		QueryParam("projection", "(id)").
		Do(&profileResp, &profileErr)

	if err != nil {
		return "", err
	}

	if !profileRes.IsSuccess() {
		return "", apperror.Internal(profileErr.Message)
	}

	postData := LinkedinPostRequest{
		Author:         "urn:li:person:" + profileResp.ID,
		LifecycleState: "PUBLISHED",
		SpecificContent: struct {
			ShareContent struct {
				ShareCommentary struct {
					Text string `json:"text"`
				} `json:"shareCommentary"`
				ShareMediaCategory string `json:"shareMediaCategory"`
			} `json:"com.linkedin.ugc.ShareContent"`
		}{
			ShareContent: struct {
				ShareCommentary struct {
					Text string `json:"text"`
				} `json:"shareCommentary"`
				ShareMediaCategory string `json:"shareMediaCategory"`
			}{
				ShareCommentary: struct {
					Text string `json:"text"`
				}{
					Text: input.Text,
				},
				ShareMediaCategory: "NONE",
			},
		},
		Visibility: struct {
			MemberNetworkVisibility string `json:"com.linkedin.ugc.MemberNetworkVisibility"`
		}{
			MemberNetworkVisibility: input.Visibility,
		},
	}

	var postResp LinkedinPostResponse
	var postErr LinkedinPostError

	postRes, err := client.Post("/v2/ugcPosts").
		JSONContentType().
		Body(postData).
		Do(&postResp, &postErr)

	if err != nil {
		return "", err
	}

	if !postRes.IsSuccess() {
		return "", apperror.Internal(postErr.Message)
	}

	return postResp.ID, nil
}
