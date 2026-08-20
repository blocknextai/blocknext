package pricemonitor

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/coingecko/helpers"
)

var (
	ErrExecutorEmptyResponse = apperror.Internal("empty response")
)

type CoingeckoPriceMonitorExecutorInput struct {
	CoinID         string  `schema:"coinId"`
	Currency       string  `schema:"currency"`
	AlertType      string  `schema:"alertType"`
	ThresholdValue float64 `schema:"thresholdValue"`
}

type CoingeckoPriceMonitorExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[CoingeckoPriceMonitorExecutorInput]
}

func NewCoingeckoPriceMonitorExecutor(
	nodeID string,
	validator *jsonschema.Validator[CoingeckoPriceMonitorExecutorInput],
) *CoingeckoPriceMonitorExecutor {
	return &CoingeckoPriceMonitorExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

type PriceResponse struct {
	Data map[string]map[string]any `json:"data"`
}

type SimplePrice map[string]map[string]float64

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}

func (e *CoingeckoPriceMonitorExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			var successResponse SimplePrice
			var errorResponse ErrorResponse
			response, err := client.Get("/simple/price").
				JSONContentType().
				QueryParam("ids", input.CoinID).
				QueryParam("vs_currencies", input.Currency).
				QueryParam("include_24hr_change", "true").
				QueryParam("include_24hr_vol", "true").
				QueryParam("include_market_cap", "true").
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			coinData, exists := successResponse[input.CoinID]
			if !exists {
				return nil, ErrExecutorEmptyResponse
			}

			currentPrice, hasPrice := coinData[input.Currency]
			if !hasPrice {
				return nil, ErrExecutorEmptyResponse
			}

			change24h := coinData[input.Currency+"_24h_change"]
			volume24h := coinData[input.Currency+"_24h_vol"]
			marketCap := coinData[input.Currency+"_market_cap"]

			var shouldAlert bool

			switch input.AlertType {
			case "price_above":
				shouldAlert = currentPrice > input.ThresholdValue
			case "price_below":
				shouldAlert = currentPrice < input.ThresholdValue
			case "change_above":
				shouldAlert = change24h > input.ThresholdValue
			case "change_below":
				shouldAlert = change24h < input.ThresholdValue
			default:
				shouldAlert = true
			}

			changeDirection := "neutral"
			if change24h > 0 {
				changeDirection = "up"
			} else if change24h < 0 {
				changeDirection = "down"
			}

			results = append(results, map[string]any{
				"coinId":          input.CoinID,
				"currency":        input.Currency,
				"currentPrice":    currentPrice,
				"change24h":       change24h,
				"volume24h":       volume24h,
				"marketCap":       marketCap,
				"alertType":       input.AlertType,
				"thresholdValue":  input.ThresholdValue,
				"shouldAlert":     shouldAlert,
				"changeDirection": changeDirection,
			})
		}

		return results, nil
	}
}
