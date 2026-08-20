package publishmediapost

import (
	"context"
	"io"
	"net/url"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/linkedin/helpers"
)

var (
	errFailedToDownloadMedia    = apperror.Internal("failed to download media")
	errFailedToReadMediaContent = apperror.Internal("failed to read media content")
)

type LinkedinPublishMediaPostExecutorInput struct {
	Text       string `schema:"text"`
	MediaURL   string `schema:"mediaUrl"`
	MediaType  string `schema:"mediaType"`
	Visibility string `schema:"visibility"`
}

type LinkedinPublishMediaPostExecutor struct {
	executors.Executor
	validator   *jsonschema.Validator[LinkedinPublishMediaPostExecutorInput]
	fileGateway filegateway.FileGateway
}

type LinkedinProfileResponse struct {
	ID string `json:"id"`
}

type LinkedinProfileError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type LinkedinMediaPostRequest struct {
	Author          string `json:"author"`
	LifecycleState  string `json:"lifecycleState"`
	SpecificContent struct {
		ShareContent struct {
			ShareCommentary struct {
				Text string `json:"text"`
			} `json:"shareCommentary"`
			ShareMediaCategory string `json:"shareMediaCategory"`
			Media              []struct {
				Status      string `json:"status"`
				Description struct {
					Text string `json:"text"`
				} `json:"description"`
				Media string `json:"media"`
				Title struct {
					Text string `json:"text"`
				} `json:"title"`
			} `json:"media"`
		} `json:"com.linkedin.ugc.ShareContent"`
	} `json:"specificContent"`
	Visibility struct {
		MemberNetworkVisibility string `json:"com.linkedin.ugc.MemberNetworkVisibility"`
	} `json:"visibility"`
}

type LinkedinMediaPostResponse struct {
	ID string `json:"id"`
}

type LinkedinMediaPostError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewLinkedinPublishMediaPostExecutor(
	nodeID string,
	validator *jsonschema.Validator[LinkedinPublishMediaPostExecutorInput],
	fileGateway filegateway.FileGateway,
) *LinkedinPublishMediaPostExecutor {
	return &LinkedinPublishMediaPostExecutor{
		ID:          nodeID,
		validator:   validator,
		fileGateway: fileGateway,
	}
}

func (e *LinkedinPublishMediaPostExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			downloadResult, err := e.fileGateway.DownloadFile(ctx, input.MediaURL)
			if err != nil {
				return nil, errFailedToDownloadMedia
			}
			defer func() {
				_ = downloadResult.DataReader.Close()
			}()

			mediaData, err := io.ReadAll(downloadResult.DataReader)
			if err != nil {
				return nil, errFailedToReadMediaContent
			}

			contentType := downloadResult.ContentType
			if contentType == "" {
				if input.MediaType == "VIDEO" {
					contentType = "video/mp4"
				} else {
					contentType = "image/jpeg"
				}
			}

			postID, err := e.uploadMediaAndCreatePost(client, *input, mediaData, contentType)
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

func (e *LinkedinPublishMediaPostExecutor) uploadMediaAndCreatePost(client *httpclient.Client, input LinkedinPublishMediaPostExecutorInput, mediaData []byte, contentType string) (string, error) {
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

	mediaCategory := "IMAGE"
	if input.MediaType == "VIDEO" {
		mediaCategory = "VIDEO"
	}

	postData := LinkedinMediaPostRequest{
		Author:         "urn:li:person:" + profileResp.ID,
		LifecycleState: "PUBLISHED",
		SpecificContent: struct {
			ShareContent struct {
				ShareCommentary struct {
					Text string `json:"text"`
				} `json:"shareCommentary"`
				ShareMediaCategory string `json:"shareMediaCategory"`
				Media              []struct {
					Status      string `json:"status"`
					Description struct {
						Text string `json:"text"`
					} `json:"description"`
					Media string `json:"media"`
					Title struct {
						Text string `json:"text"`
					} `json:"title"`
				} `json:"media"`
			} `json:"com.linkedin.ugc.ShareContent"`
		}{
			ShareContent: struct {
				ShareCommentary struct {
					Text string `json:"text"`
				} `json:"shareCommentary"`
				ShareMediaCategory string `json:"shareMediaCategory"`
				Media              []struct {
					Status      string `json:"status"`
					Description struct {
						Text string `json:"text"`
					} `json:"description"`
					Media string `json:"media"`
					Title struct {
						Text string `json:"text"`
					} `json:"title"`
				} `json:"media"`
			}{
				ShareCommentary: struct {
					Text string `json:"text"`
				}{
					Text: input.Text,
				},
				ShareMediaCategory: mediaCategory,
				Media: []struct {
					Status      string `json:"status"`
					Description struct {
						Text string `json:"text"`
					} `json:"description"`
					Media string `json:"media"`
					Title struct {
						Text string `json:"text"`
					} `json:"title"`
				}{
					{
						Status: "READY",
						Description: struct {
							Text string `json:"text"`
						}{
							Text: input.Text,
						},
						Media: input.MediaURL,
						Title: struct {
							Text string `json:"text"`
						}{
							Text: "Media Post",
						},
					},
				},
			},
		},
		Visibility: struct {
			MemberNetworkVisibility string `json:"com.linkedin.ugc.MemberNetworkVisibility"`
		}{
			MemberNetworkVisibility: input.Visibility,
		},
	}

	var postResp LinkedinMediaPostResponse
	var postErr LinkedinMediaPostError

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
