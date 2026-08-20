package veo

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/base64"
	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/veo/helpers"
)

var (
	ErrFailedToDownloadVideo = apperror.Internal("failed to download video")
)

type VeoExecutorInput struct {
	Model       string `schema:"model"`
	Prompt      string `schema:"prompt"`
	Image       string `schema:"image"`
	AspectRatio string `schema:"aspectRatio"`
	//PersonGeneration string `schema:"personGeneration"`
}

type VeoExecutor struct {
	executors.Executor
	validator   *jsonschema.Validator[VeoExecutorInput]
	fileGateway filegateway.FileGateway
}

func NewVeoExecutor(
	nodeID string,
	validator *jsonschema.Validator[VeoExecutorInput],
	fileGateway filegateway.FileGateway,
) *VeoExecutor {
	return &VeoExecutor{
		ID:          nodeID,
		validator:   validator,
		fileGateway: fileGateway,
	}
}

type VeoSuccessResponse struct {
	Name     string `json:"name"`
	Metadata struct {
		Type        string `json:"@type"`
		CreateTime  string `json:"createTime"`
		UpdateTime  string `json:"updateTime"`
		ProgressVal int    `json:"progressPercentage"`
	} `json:"metadata"`
	Done     bool `json:"done"`
	Response struct {
		Type                  string `json:"@type"`
		GenerateVideoResponse struct {
			GeneratedSamples []struct {
				Video struct {
					URI string `json:"uri"`
				} `json:"video"`
			} `json:"generatedSamples"`
		} `json:"generateVideoResponse"`
	} `json:"response"`
}

type VeoErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (e *VeoExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "veo_api")
		apiKey := credential.String("apiKey")
		client := helpers.CreateClient(ctx, apiKey)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			instance := map[string]any{}

			if strings.TrimSpace(input.Prompt) != "" {
				instance["prompt"] = input.Prompt
			}

			if strings.TrimSpace(input.Image) != "" {
				imageFile, err := e.fileGateway.DownloadFile(ctx, input.Image)
				if err != nil {
					return nil, err
				}
				defer func() {
					_ = imageFile.DataReader.Close()
				}()

				base64Data, err := base64.EncodeReader(imageFile.DataReader)
				if err != nil {
					return nil, err
				}

				instance["image"] = map[string]any{
					"bytesBase64Encoded": base64Data,
					"mimeType":           imageFile.ContentType,
				}
			}

			var successResponse VeoSuccessResponse
			var errorResponse VeoErrorResponse
			response, err := client.Post("/models/"+input.Model+":predictLongRunning").
				JSONContentType().
				Body(map[string]any{
					"instances": []map[string]any{instance},
					"parameters": map[string]any{
						"aspectRatio": input.AspectRatio,
						//"personGeneration": input.PersonGeneration,
					},
				}).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() || successResponse.Name == "" {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			videoURLs, err := WaitVeo(ctx, client, WaitVeoInput{
				OperationName: successResponse.Name,
				APIKey:        apiKey,
			})
			if err != nil {
				return nil, err
			}

			for index, videoURL := range videoURLs {
				downloadedURL, err := e.downloadAndUploadVideo(ctx, videoURL, index)
				if err != nil {
					return nil, err
				}

				results = append(results, map[string]any{
					"video": downloadedURL,
				})
			}
		}

		return results, nil
	}
}

func (e *VeoExecutor) downloadAndUploadVideo(ctx context.Context, videoURL string, index int) (string, error) {
	client := httpclient.NewClientBuilder().
		Context(ctx).
		Build()

	response, err := client.Get(videoURL).DoStream()
	if err != nil {
		return "", err
	}
	defer func() {
		err := response.BodyReader.Close()
		if err != nil {
			slog.Error("Failed to close response body reader",
				"component", "VeoExecutor",
				"error", err)
		}
	}()

	if !response.IsSuccess() {
		return "", ErrFailedToDownloadVideo
	}

	uploadResult, err := e.fileGateway.UploadFile(
		ctx,
		"60b9c9f6-e2e6-4782-92a0-8ba2b8c32720",
		"veo-"+strconv.Itoa(index)+".mp4",
		response.BodyReader,
	)
	if err != nil {
		return "", err
	}

	return uploadResult.URL, nil
}
