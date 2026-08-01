package airdroptracker

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/coingecko/helpers"
)

var (
	ErrExecutorEmptyResponse = apperror.Internal("empty response")
)

type CoingeckoAirdropTrackerExecutorInput struct {
	ProjectName  string  `schema:"projectName"`
	StatusFilter string  `schema:"statusFilter"`
	MinReward    float64 `schema:"minReward"`
}

type CoingeckoAirdropTrackerExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[CoingeckoAirdropTrackerExecutorInput]
}

func NewCoingeckoAirdropTrackerExecutor(
	nodeID string,
	validator *jsonschema.Validator[CoingeckoAirdropTrackerExecutorInput],
) *CoingeckoAirdropTrackerExecutor {
	return &CoingeckoAirdropTrackerExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type AirdropResponse struct {
	Data []struct {
		ID              string    `json:"id"`
		Name            string    `json:"name"`
		Description     string    `json:"description"`
		Status          string    `json:"status"`
		StartDate       time.Time `json:"start_date"`
		EndDate         time.Time `json:"end_date"`
		ProjectName     string    `json:"project_name"`
		ProjectWebsite  string    `json:"project_website"`
		ProjectTwitter  string    `json:"project_twitter"`
		ProjectTelegram string    `json:"project_telegram"`
		ProjectDiscord  string    `json:"project_discord"`
		RewardAmount    string    `json:"reward_amount"`
		RewardCurrency  string    `json:"reward_currency"`
		Requirements    []string  `json:"requirements"`
		Instructions    string    `json:"instructions"`
	} `json:"data"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}

func (e *CoingeckoAirdropTrackerExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "coingecko_api")
		apiKey := credential.String("apiKey")
		client := helpers.CreateClient(ctx, apiKey)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			var successResponse AirdropResponse
			var errorResponse ErrorResponse
			response, err := client.Get("/airdrops").
				JSONContentType().
				QueryParam("status", input.StatusFilter).
				QueryParam("per_page", "50").
				QueryParam("page", "1").
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			if len(successResponse.Data) == 0 {
				return nil, ErrExecutorEmptyResponse
			}

			var filteredAirdrops []map[string]any
			for _, airdrop := range successResponse.Data {
				if input.ProjectName != "" && !strings.Contains(strings.ToLower(airdrop.ProjectName), strings.ToLower(input.ProjectName)) {
					continue
				}

				rewardValue := 0.0
				if airdrop.RewardAmount != "" {
					if val, err := strconv.ParseFloat(strings.TrimSpace(airdrop.RewardAmount), 64); err == nil {
						rewardValue = val
					}
				}

				if rewardValue < input.MinReward {
					continue
				}

				airdropData := map[string]any{
					"id":              airdrop.ID,
					"name":            airdrop.Name,
					"description":     airdrop.Description,
					"status":          airdrop.Status,
					"startDate":       airdrop.StartDate.Format("2006-01-02"),
					"endDate":         airdrop.EndDate.Format("2006-01-02"),
					"projectName":     airdrop.ProjectName,
					"projectWebsite":  airdrop.ProjectWebsite,
					"projectTwitter":  airdrop.ProjectTwitter,
					"projectTelegram": airdrop.ProjectTelegram,
					"projectDiscord":  airdrop.ProjectDiscord,
					"rewardAmount":    airdrop.RewardAmount,
					"rewardCurrency":  airdrop.RewardCurrency,
					"rewardValue":     rewardValue,
					"requirements":    airdrop.Requirements,
					"instructions":    airdrop.Instructions,
				}

				filteredAirdrops = append(filteredAirdrops, airdropData)
			}

			results = append(results, map[string]any{
				"airdrops":      filteredAirdrops,
				"totalCount":    len(filteredAirdrops),
				"projectFilter": input.ProjectName,
				"statusFilter":  input.StatusFilter,
				"minReward":     input.MinReward,
			})
		}

		return results, nil
	}
}
