package createplaylist

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/soundcloud/helpers"
)

var (
	ErrExecutorEmptyResponse = apperror.Internal("empty response")
	ErrCreateFailed          = apperror.Internal("create failed")
)

type SoundCloudCreatePlaylistExecutorInput struct {
	Title       string `schema:"title"`
	Description string `schema:"description"`
	TrackIDs    string `schema:"trackIds"`
}

type SoundCloudCreatePlaylistExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[SoundCloudCreatePlaylistExecutorInput]
}

func NewSoundCloudCreatePlaylistExecutor(
	nodeID string,
	validator *jsonschema.Validator[SoundCloudCreatePlaylistExecutorInput],
) *SoundCloudCreatePlaylistExecutor {
	return &SoundCloudCreatePlaylistExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

func (e *SoundCloudCreatePlaylistExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "soundcloud_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			playlistResponse, err := e.createPlaylist(client, *input)
			if err != nil {
				return nil, err
			}

			if playlistResponse == nil {
				return nil, ErrExecutorEmptyResponse
			}

			results = append(results, map[string]any{
				"status": true,
			})
		}

		return results, nil
	}
}

type PlaylistResponse struct {
	ID           int    `json:"id"`
	PermalinkURL string `json:"permalink_url"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}

type ErrorResponse struct {
	Errors  []map[string]any `json:"errors"`
	Message string           `json:"message"`
}

func (e *SoundCloudCreatePlaylistExecutor) createPlaylist(client *httpclient.Client, input SoundCloudCreatePlaylistExecutorInput) (*PlaylistResponse, error) {
	tracks := make([]map[string]any, 0)
	if len(input.TrackIDs) > 0 {
		for _, id := range input.TrackIDs {
			tracks = append(tracks, map[string]any{
				"id": id,
			})
		}
	}

	playlistPayload := make(map[string]any)
	playlistPayload["title"] = input.Title
	if input.Description != "" {
		playlistPayload["description"] = input.Description
	}
	if len(tracks) > 0 {
		playlistPayload["tracks"] = tracks
	}
	payload := map[string]any{
		"playlist": playlistPayload,
	}

	var successResponse PlaylistResponse
	var errorResponse ErrorResponse
	response, err := client.Post("/playlists").
		JSONContentType().
		Body(payload).
		Do(&successResponse, &errorResponse)

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, ErrCreateFailed
	}

	return &successResponse, nil
}
